package adapter

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"

	pluginApi "github.com/MaibornWolff/elcep/plugin"
)

// PluginProvider loads the plugin files and scans for available plugins
type PluginProvider struct {
	Plugins map[string]func(interface{}) pluginApi.Plugin
}

// NewPluginProvider returns an instance with loaded Plugins from plugin Files
func NewPluginProvider(pluginFolder string) *PluginProvider {
	provider := &PluginProvider{}
	files := findPlugins(pluginFolder)
	provider.initializePlugins(files)
	return provider
}

// GetPluginNames returns a list of logical plugin names
func (provider *PluginProvider) GetPluginNames() []string {
	keys := make([]string, 0, len(provider.Plugins))
	for k := range provider.Plugins {
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
	provider.Plugins = make(map[string]func(interface{}) pluginApi.Plugin)
	for _, file := range fileNames {
		plug, err := plugin.Open(file)
		if err != nil {
			log.Fatalf("%s: os.Open(): %s\n", file, err)
		}

		sym, err := plug.Lookup("NewPlugin")
		if err != nil {
			log.Fatal(err)
		}

		m, ok := sym.(func(interface{}) pluginApi.Plugin)
		if !ok {
			var expected func(interface{}) pluginApi.Plugin
			log.Fatalf("unexpected type from module symbol NewPlugin. Expected `%T`", expected)
		}

		pluginName := getLogicalPluginName(file)
		provider.Plugins[pluginName] = m
	}
}

func getLogicalPluginName(file string) string {
	name := filepath.Base(file)
	return name[0 : len(name)-3]
}
