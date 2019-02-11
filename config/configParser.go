package config

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

type configurationFile struct {
	Plugins map[string]interface{}
	Metrics map[string]queryGroup
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
	err := yaml.Unmarshal(fileContent, &configFile)
	if err != nil {
		log.Fatalln("Could not parse config file")
	}
	fmt.Printf("--- Loaded yaml:\n%#v\n\n", configFile)

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
				} else if queryMap, ok := query.(Query); ok {
					if _, ok := queryMap["Name"]; !ok {
						queryMap["Name"] = name
					}
					if !queryMap.isValid() {
						log.Fatalln("Plugin Config for plugin ", pluginName, " is not valid, missing 'Query' or 'Name' of type string")
					}
					pluginConf.Queries = append(pluginConf.Queries, queryMap)
				} else {
					log.Fatalln("Plugin Config for plugin ", pluginName, " is of wrong type")
				}
			}
		}
	}
	return
}
