package adapter

import (
	"flag"
	"log"
)

//CommandLineOption defines possible options
type CommandLineOption struct {
	Freq                   int
	ElasticsearchURL       string
	Port                   int
	Path                   string
	ElasticsearchQueryFile string
	TimeKey                string
}

func defaultOptions() *CommandLineOption {
	return &CommandLineOption{
		Freq:                   30,
		ElasticsearchURL:       "http://elasticsearch:9200",
		Port:                   8080,
		Path:                   "/metrics",
		ElasticsearchQueryFile: "./conf/queries.cfg",
		TimeKey:                "@timestamp",
	}
}

// ParseOptions returns the default options, overriden by the command line options if present
func ParseOptions() *CommandLineOption {
	var options = defaultOptions()

	flag.StringVar(&options.ElasticsearchURL, "url", options.ElasticsearchURL, "The elastic search endpoint")
	flag.IntVar(&options.Port, "port", options.Port, "The port to listen on for HTTP requests")
	flag.StringVar(&options.Path, "path", options.Path, "The path to listen on for HTTP requests")
	flag.IntVar(&options.Freq, "freq", options.Freq, "The interval in seconds in which to query elastic search")
	flag.StringVar(&options.ElasticsearchQueryFile, "elasticqueries", options.ElasticsearchQueryFile, "The path to the queries.cfg")
	flag.StringVar(&options.TimeKey, "time-key", options.TimeKey, "The time key to use in elastic search queries")
	flag.Parse()

	return options
}

//PrintCmdLineOptions as logs
func (option *CommandLineOption) PrintCmdLineOptions() {
	log.Println("Config:")
	log.Println("\tUrl:", option.ElasticsearchURL)
	log.Println("\tFreq:", option.Freq)
	log.Println("\tPort:", option.Port)
	log.Println("\tElasticsearch Configuration File:", option.ElasticsearchQueryFile)
	log.Println("\tTime Key:", option.TimeKey)
}
