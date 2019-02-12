package config

import (
	"log"

	"gopkg.in/yaml.v2"
)

type configurationFile struct {
	Plugins map[string]interface{} `yaml:"plugins"`
	Metrics map[string]queryGroup  `yaml:"metrics"`
}

type queryGroup map[string]queries

type queries map[string]interface{}

// PluginConfig the config struct for a plugin
type PluginConfig struct {
	Options interface{}
	Queries []Query
}

func parseConfigFile(fileContent []byte) (conf map[string]*PluginConfig) {
	var configFile configurationFile
	err := yaml.UnmarshalStrict(fileContent, &configFile)
	if err != nil {
		log.Fatalf("Could not parse config file: %v\n", err)
	}

	conf = make(map[string]*PluginConfig)
	for pluginName, pluginConf := range configFile.Plugins {
		conf[pluginName] = &PluginConfig{}
		conf[pluginName].Options = pluginConf
	}

	for groupName, group := range configFile.Metrics {
		for pluginName, queries := range group {
			pluginConf := conf[pluginName]
			if pluginConf == nil {
				pluginConf = &PluginConfig{}
			}
			for queryName, query := range queries {
				name := groupName + "_" + queryName
				if queryString, ok := query.(string); ok {
					q := CreateQuery(name, queryString)
					if !q.isValid() {
						log.Fatalf("Query is invalid: %#v\n", q)
					}
					pluginConf.Queries = append(pluginConf.Queries, q)
				} else if queryMap, ok := query.(map[interface{}]interface{}); ok {
					if qObj, ok := queryMap["query"]; !ok {
						log.Fatalln("Plugin Config for plugin ", pluginName, " is not valid, missing 'Query' for ", name)
					} else if qString, ok := qObj.(string); !ok {
						log.Fatalln("Plugin Config for plugin ", pluginName, " is not valid, 'Query' must be a string.")
					} else {
						q := CreateQuery(name, qString)
						pluginConf.Queries = append(pluginConf.Queries, q)
					}
				} else {
					log.Fatalln("Plugin Config for plugin ", pluginName, " is of wrong type")
				}
			}
		}
	}
	return
}
