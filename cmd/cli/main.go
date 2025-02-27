package main

import (
	"achaean/plugins"
	"fmt"
	"os"
	"path/filepath"
)

const PLUGINS_PATH = "./data/plugins"

func main2() {
	// New PluginManager using the plugins path
	pm := plugins.NewPluginManager()
	pluginsPath, err := filepath.Abs(PLUGINS_PATH)
	if err != nil {
		println(err.Error())
		return
	}

	// Exit if error loading plugins
	err = pm.LoadPluginsFromDir(pluginsPath)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	} else {
		for _, p := range pm.Plugins {
			fmt.Println(p)
		}
	}

	// If no param, show usage and exit
	if len(os.Args) == 1 {
		fmt.Printf("USAGE: %s <PluginID>\n", os.Args[0])
		return
	}

    pluginID := os.Args[1]

	// Set the functions to manage pipes output
	pm.StdoutFunc = StdoutHandle
	pm.StderrFunc = StderrHandle
	pm.ProgressPipeFunc = ProgressHandle

	// Try to run the plugin
	fmt.Println("------<BEGIN ExecPlugin>------")
	err = pm.ExecPluginByID(pluginID)
	//pm.SoftKillCurrentPlugin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Wait until plugin is finished
	for {
		if !pm.PluginIsRunning() {
			break
		}
	}
	fmt.Println("------<END ExecPlugin>------")

}

// Handle functions for the pipes.
func StdoutHandle(s string) {
	fmt.Printf("OUT: %s\n", s)
}

func StderrHandle(s string) {
	fmt.Printf("ERROR: %s\n", s)
}

func ProgressHandle(p int) {
	fmt.Printf("[%v]", p)
}
