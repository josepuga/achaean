package main

import (
	"achaean/plugins"
	"fmt"
	"path/filepath"
    
)

const PLUGINS_PATH = "./data/plugins"

func main() {
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

	// GUI App
	a := NewApp(pm)
    a.Start()


}
