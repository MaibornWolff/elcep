package adapter

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

//Configuration the configuration for the whole application
type Configuration struct {
	getReader func(string) io.ReadCloser
	byFile    map[string]PluginConfig
}

//PluginConfig the configuration for a certain plugin
type PluginConfig map[string]string

func (config *Configuration) ForPlugin(pluginName string) PluginConfig {
	return config.byFile[pluginName]
}

//Print all queries as logs
func (config *Configuration) Print() {
	log.Println("ElasticSearch Queries:")
	for file, queries := range config.byFile {
		log.Println("\t", file)
		for query := range queries {
			log.Println("\t\t", query)
		}
	}
}

//ReadConfig takes a list of plugin names and a function to obtain a Scanner for each plugin providing the configuration
func ReadConfig(pluginNames []string, configReader func(string) io.ReadCloser) *Configuration {
	config := &Configuration{
		getReader: configReader,
		byFile:    make(map[string]PluginConfig),
	}
	config.read(pluginNames)
	return config
}

func (config *Configuration) read(pluginNames []string) {
	for _, plugin := range pluginNames {
		config.byFile[plugin] = config.readFromFile(plugin)
	}
}

//Read all queries from a given file
func (config *Configuration) readFromFile(pluginName string) map[string]string {
	var queries = make(map[string]string)

	re, err := regexp.Compile(`([^=]+)=(.*)`)
	if err != nil {
		log.Fatalf("%s: regexp.Compile(): error=%s", "ReadQueriesConfig()", err)
	}

	reader := config.getReader(pluginName)
	defer reader.Close()
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") && strings.Contains(line, "=") {
			queryKey := re.FindStringSubmatch(line)[1]
			queryValue := re.FindStringSubmatch(line)[2]
			queries[queryKey] = queryValue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("%s: scanner.Err(): %s\n", "ReadQueriesConfig", err)
		os.Exit(1)
	}

	return queries
}
