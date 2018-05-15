package adapter

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"plugin"

	"github.com/MaibornWolff/elcep/monitor"
)

type PluginProvider struct {
	pluginFolder string
	Monitors     map[string]func() monitor.LogMonitor
}

//NewPluginProvider returns an instance with loaded LogMonitors from plugin Files
func NewPluginProvider(pluginFolder string) *PluginProvider {
	provider := &PluginProvider{}
	provider.pluginFolder = pluginFolder
	files := provider.findPluginFileNames()
	provider.initializePlugins(files)
	return provider
}

//GetPluginNames returns a list of logical plugin names
func (provider *PluginProvider) GetPluginNames() []string {
	keys := make([]string, 0, len(provider.Monitors))
	for k := range provider.Monitors {
		keys = append(keys, k)
	}
	return keys
}

func (provider *PluginProvider) initializePlugins(fileNames []string) {
	provider.Monitors = make(map[string]func() monitor.LogMonitor)
	for _, file := range fileNames {
		plug, err := plugin.Open(file)
		if err != nil {
			log.Printf("%s: os.Open(): %s\n", file, err)
			os.Exit(1)
		}

		sym, err := plug.Lookup("NewMonitor")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		m, ok := sym.(func() monitor.LogMonitor)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}

		provider.Monitors[getLogicalPluginName(file)] = m
	}
}

func (provider *PluginProvider) findPluginFileNames() []string {
	var foundFileNames []string

	files, err := ioutil.ReadDir(provider.pluginFolder)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		foundFileNames = append(foundFileNames, filepath.Join(provider.pluginFolder, f.Name()))
	}

	return foundFileNames
}

func getLogicalPluginName(file string) string {
	name := filepath.Base(file)
	return name[0 : len(name)-3]
}
