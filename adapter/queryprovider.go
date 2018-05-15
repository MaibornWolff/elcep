package adapter

import (
	"bufio"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

//QueryProvider loads queries
type QueryProvider struct {
	QuerySets     map[string]QuerySet
	openQueryFile func(string) io.ReadCloser
}

//QuerySet defines a logical set for queries with a given name
type QuerySet struct {
	Queries map[string]string
}

//Print all queries as logs
func (provider *QueryProvider) Print() {
	log.Println("ElasticSearch Queries:")
	for key, set := range provider.QuerySets {
		log.Println("\t", key)
		for k := range set.Queries {
			log.Println("\t\t", k)
		}
	}
}

//NewQueryProvider takes a filename for the default Counter Monitor and a list of plugin names
func NewQueryProvider(counterQueryFile string, pluginNames []string) *QueryProvider {
	provider := &QueryProvider{
		openQueryFile: openQueryFile,
	}

	provider.read(counterQueryFile, pluginNames)

	return provider
}

func (provider *QueryProvider) read(counterQueryFile string, pluginNames []string) {
	provider.QuerySets = make(map[string]QuerySet)
	provider.QuerySets["default"] = QuerySet{
		Queries: provider.readFromFile(counterQueryFile),
	}

	for _, plugin := range pluginNames {
		provider.QuerySets[plugin] = QuerySet{
			Queries: provider.readFromFile("./conf/" + plugin + ".cfg"),
		}
	}
}

//Read all queries from a given file
func (provider *QueryProvider) readFromFile(queriesFile string) map[string]string {

	var queries map[string]string
	queries = make(map[string]string)

	file := provider.openQueryFile(queriesFile)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") == false {
			if strings.Contains(line, "=") == true {
				re, err := regexp.Compile(`([^=]+)=(.*)`)
				if err != nil {
					log.Printf("%s: regexp.Compile(): error=%s", "ReadQueriesConfig()", err)
				} else {
					queryKey := re.FindStringSubmatch(line)[1]
					queryValue := re.FindStringSubmatch(line)[2]
					queries[queryKey] = queryValue
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("%s: scanner.Err(): %s\n", "ReadQueriesConfig", err)
		os.Exit(1)
	}

	return queries
}

func openQueryFile(name string) io.ReadCloser {
	file, err := os.Open(name)
	if err != nil {
		log.Printf("%s: os.Open(): %s\n", name, err)
		os.Exit(1)
	}

	return file
}
