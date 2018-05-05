package adapter

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

const queryTemplate = `{
    "query": {
        "query_string": {
            "query": "<query>"
        }
    },
    "size":0
}`

//QueryProvider loads queries
type QueryProvider struct {
	Queries *map[string]string
}

//Print as logs
func (provider *QueryProvider) Print() {
	log.Println("ElasticSearch Queries:")
	for k := range *provider.Queries {
		log.Println("\t", k)
	}
}

//Read all queries from a given file
func (provider *QueryProvider) Read(queriesFile string) {
	prg := "ReadQueriesConfig()"

	var queries map[string]string
	queries = make(map[string]string)

	file, err := os.Open(queriesFile)
	defer file.Close()

	if err != nil {
		log.Printf("%s: os.Open(): %s\n", prg, err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") == false {
			if strings.Contains(line, "=") == true {
				re, err := regexp.Compile(`([^=]+)=(.*)`)
				if err != nil {
					log.Printf("%s: regexp.Compile(): error=%s", prg, err)
				} else {
					queryKey := re.FindStringSubmatch(line)[1]
					queryValue := re.FindStringSubmatch(line)[2]
					queries[queryKey] = provider.buildFromTemplate(queryValue)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("%s: scanner.Err(): %s\n", prg, err)
	}

	provider.Queries = &queries
}

func (provider *QueryProvider) buildFromTemplate(query string) string {
	return strings.Replace(queryTemplate, "<query>", query, 1)
}
