package plugins

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type PluginManager struct {
	Plugins                map[string]*Plugin
	currentPlugin          *Plugin   // Current running plugin
	command                *exec.Cmd // Current running plugin command
	stdout, stderr         io.ReadCloser
	StdoutFunc, StderrFunc func(string)
	ProgressPipe           *os.File
	ProgressPipeFunc       func(p int) // 0 - 100
	pluginIsRunning        bool
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: make(map[string]*Plugin),
	}
}

func (pm *PluginManager) GetPluginByID(id string) (*Plugin, error) {
    var e error
    p, ok := pm.Plugins[id]
    if !ok {
        e = fmt.Errorf("plugin id %s does not exists", id) 
    }
    return p, e
}

func (pm *PluginManager) ExecPlugin( p *Plugin) error {
    return pm.ExecPluginByID(p.ID)
}

func (pm *PluginManager) ExecPluginByID(id string) error {
	p, err := pm.GetPluginByID(id)
	if err != nil {
		return err
	}
	pm.currentPlugin = p
	pm.command = exec.Command(p.Entrypoint, p.GetParametersAsSlice()...)
	pm.command.Dir = pm.currentPlugin.Dir
	pm.command.Env = append(os.Environ(),
		"PYTHONUNBUFFERED=1",  // Sólo válido para python que usa su propio buffer
		fmt.Sprintf("PLUGIN_ID=%s", p.ID),
		fmt.Sprintf("PROGRESS_PIPE=%s", GetProgressPipeName(p.ID)),
	)

	// Crea un grupo de procesos para manejar el plugin y sus hijos
	pm.command.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return pm.startCurrentPlugin()
}

func (pm *PluginManager) PluginIsRunning() bool {
	return pm.pluginIsRunning
}

func (pm *PluginManager) handleRunningPlugin() {
	var wg sync.WaitGroup
	stopChan := make(chan struct{}) // Trigger para detener las goroutines locales
	defer pm.closePipes()

	// Goroutine para leer `progressPipe`
	handleProgress := func() {
		defer wg.Done()
		buf := make([]byte, 8)
		for {
			select {
			case <-stopChan:
				return
			default:
				n, err := pm.ProgressPipe.Read(buf)
				if err != nil {
					if err != io.EOF {
						fmt.Fprintf(os.Stderr, "error reading ProgressPipe: %s\n", err.Error())
					}
					return
				}
				if n > 0 {
					progress := strings.TrimSpace(string(buf[:n]))
					if progress == "DONE" {
						close(stopChan)
						return
					}
					if num, err := strconv.Atoi(progress); err == nil {
						if pm.ProgressPipeFunc != nil {
							pm.ProgressPipeFunc(num)
						}
					} else {
						fmt.Fprintf(os.Stderr, "error converting progress number: %s\n", err.Error())
					}
				}
			}
		}
	}

	// Goroutines for `stdout` & `stderr`
	handleOutput := func(reader io.Reader, stdHandleFunc func(string)) {
		defer wg.Done()
		scanner := bufio.NewScanner(reader) // Read line by line
		for scanner.Scan() {
			if stdHandleFunc != nil {
				stdHandleFunc(scanner.Text())
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading output: %s\n", err.Error())
		}
	}

	wg.Add(3)
	go handleProgress()
	go handleOutput(pm.stdout, pm.StdoutFunc)
	go handleOutput(pm.stderr, pm.StderrFunc)

	wg.Wait()                                   // Espera a que terminen las goroutines
	pm.pluginIsRunning = false
}

func (pm *PluginManager) startCurrentPlugin() error {
	err := pm.openPipes()
	if err != nil {
		pm.closePipes()
		return err
	}
	err = pm.command.Start()
	if err != nil {
		pm.closePipes()
		return err
	}
	pm.pluginIsRunning = true
	go pm.handleRunningPlugin()
	return nil
}

// SIGTERM (15)
func (pm *PluginManager) SoftKillCurrentPlugin() error {
	// Send SIGTERM to the process group. Minus sign "-" is for send the signal to all group.
	return syscall.Kill(-pm.command.Process.Pid, syscall.SIGTERM)
}

// SIGKILL (9)
func (pm *PluginManager) HardKillCurrentPlugin() error {
	// Send SIGKILL to the process group. Minus sign "-" is for send the signal to all group.
	return syscall.Kill(-pm.command.Process.Pid, syscall.SIGKILL)
}

func (pm *PluginManager) openPipes() error {
	var err error

	// Progress Pipe
	pm.ProgressPipe, err = pm.getProgressPipe()
	if err != nil {
		return err
	}

	// Capture stdout from plugin
	pm.stdout, err = pm.command.StdoutPipe()
	if err != nil {
		return err
	}

	// Capture stderr from plugin
	pm.stderr, err = pm.command.StderrPipe()
	if err != nil {
		return err
	}
	return nil
}

func (pm *PluginManager) closePipes() {
	if pm.ProgressPipe != nil {
		pm.ProgressPipe.Close()
	}
	if pm.stdout != nil {
		pm.stdout.Close()
	}
	if pm.stderr != nil {
		pm.stderr.Close()
	}
}

func (pm *PluginManager) getProgressPipe() (*os.File, error) {
	var pipeFile = GetProgressPipeName(pm.currentPlugin.ID)

	// Delete the (possible) old pipe
	if err := os.Remove(pipeFile); err != nil && !os.IsNotExist(err) {
		// If the file exists and cannot be removed, send error
		return nil, err
	}

	// Create the named pipe
	err := syscall.Mkfifo(pipeFile, 0666)
	if err != nil {
		return nil, err
	}

	// Open the pipe in RW mode (RO halt the program waiting for input)
    // os.O_SYNC = no buffer.
	progressPipe, err := os.OpenFile(pipeFile, os.O_RDWR|os.O_SYNC, os.ModeNamedPipe)
	if err != nil {
		return nil, err
	}

	return progressPipe, nil
}

func (pm *PluginManager) Register(plugin *Plugin) {
	pm.Plugins[plugin.ID] = plugin
}

func (pm *PluginManager) GetCategories() []string {
	mapCat := make(map[string]interface{})
	// Using map for unique value
	for _, p := range pm.Plugins {
		mapCat[p.Category] = nil
	}

	// Convert map to []string
	result := []string{}
	for cat := range mapCat {
		result = append(result, cat)
	}
	return result
}

func (pm *PluginManager) GetPluginsByCategory(cat string) []*Plugin {
	result := []*Plugin{}
	for _, p := range pm.Plugins {
		if p.Category == cat {
			result = append(result, p)
		}
	}
	return result
}

func (pm *PluginManager) JsonToPlugin(jsonContent []byte) *Plugin {
	result := NewPlugin()
	//TODO: Coger los valores del JSON
	return result
}

func (pm *PluginManager) LoadPluginsFromDir(path string) error {

	//Category names: Only First Uppercase. Only numbers, letters and "-" allowed
	var categoryRegex = regexp.MustCompile(`^[A-Z][A-Za-z0-9 -]*$`)

	// Get directories inside plugins dir
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			pluginDir := filepath.Join(path, entry.Name())
			jsonFile := filepath.Join(pluginDir, "plugin.json")
			// Get the JSON content in a []byte
			jsonContent, err := os.ReadFile(jsonFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to read plugin json %s: %s\n", pluginDir, err.Error())
				continue
			}
			p, err := NewPluginFromJson(jsonContent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot create plugin from json at %s: %s\n", pluginDir, err.Error())
				continue
			}
			if p.ID == "" {
				fmt.Fprintf(os.Stderr, "plugin at %s has an empty ID\n", pluginDir)
				continue
			}
			if p.Category == "" {
				fmt.Fprintf(os.Stderr, "plugin at %s has an empty category\n", pluginDir)
				continue
			}
			if !categoryRegex.MatchString(p.Category) {
				fmt.Fprintf(os.Stderr, "invalid category name format in plugin at %s: %s\n", pluginDir, p.Category)
				continue
			}
			p.Dir = filepath.Join(path, entry.Name())
			p.Entrypoint = filepath.Join(p.Dir, p.Entrypoint)
			pm.Register(p)
		}
	}
	if len(pm.Plugins) == 0 {
		return fmt.Errorf("%s does not contains plugins", path)
	}
	return nil
}
