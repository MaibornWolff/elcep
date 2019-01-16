package adapter

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"

	"github.com/MaibornWolff/elcep/monitor"
)

type PluginProvider struct {
	Monitors map[string]func() monitor.LogMonitor
}

//NewPluginProvider returns an instance with loaded LogMonitors from plugin Files
func NewPluginProvider(pluginFolder string) *PluginProvider {
	provider := &PluginProvider{}
	files := findPlugins(pluginFolder)
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

func findPlugins(pluginFolder string) []string {
	var foundFileNames []string

	if files, err := ioutil.ReadDir(pluginFolder); err != nil {
		log.Fatal(err)
	} else {
		for _, f := range files {
			foundFileNames = append(foundFileNames, filepath.Join(pluginFolder, f.Name()))
		}
	}

	return foundFileNames
}

func (provider *PluginProvider) initializePlugins(fileNames []string) {
	provider.Monitors = make(map[string]func() monitor.LogMonitor)
	for _, file := range fileNames {
		plug, err := plugin.Open(file)
		if err != nil {
			log.Fatalf("%s: os.Open(): %s\n", file, err)
		}

		sym, err := plug.Lookup("NewMonitor")
		if err != nil {
			log.Fatal(err)
		}

		m, ok := sym.(func() monitor.LogMonitor)
		if !ok {
			log.Fatal("unexpected type from module symbol NewMonitor. Expected `monitor.LogMonitor`")
		}

		pluginName := getLogicalPluginName(file)
		provider.Monitors[pluginName] = m
	}
}

func getLogicalPluginName(file string) string {
	name := filepath.Base(file)
	return name[0 : len(name)-3]
}
