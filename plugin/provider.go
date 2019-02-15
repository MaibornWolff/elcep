package plugin

import (
	"github.com/MaibornWolff/elcep/config"
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"
)

const factoryMethodName = "NewPlugin"

// this is simply a type alias. This is not a new type!
type factoryMethodType = func(config.Options, interface{}) Plugin

// Provider loads the plugin files and scans for available plugins
type Provider struct {
	Plugins map[string]factoryMethodType
}

// NewPluginProvider returns an instance with loaded Plugins from plugin Files
func NewPluginProvider(pluginFolder string) *Provider {
	provider := &Provider{}
	files := findPlugins(pluginFolder)
	provider.initializePlugins(files)
	return provider
}

// GetPluginNames returns a list of logical plugin names
func (provider *Provider) GetPluginNames() []string {
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

func (provider *Provider) initializePlugins(fileNames []string) {
	provider.Plugins = make(map[string]factoryMethodType)
	for _, file := range fileNames {
		plug, err := plugin.Open(file)
		if err != nil {
			log.Fatalf("%s: os.Open(): %s\n", file, err)
		}

		sym, err := plug.Lookup(factoryMethodName)
		if err != nil {
			log.Fatalf("%s: Could not find symbol '%s': %s\n", file, factoryMethodName, err)
		}

		m, ok := sym.(factoryMethodType)
		if !ok {
			var expected factoryMethodType
			log.Fatalf("%s: unexpected type from module symbol %s. Expected `%T`", file, factoryMethodName, expected)
		}

		pluginName := getLogicalPluginName(file)
		provider.Plugins[pluginName] = m
	}
}

func getLogicalPluginName(file string) string {
	name := filepath.Base(file)
	return name[0 : len(name)-3]
}
