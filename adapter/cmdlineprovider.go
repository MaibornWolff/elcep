package adapter

import (
	"flag"
	"log"
)

//CommandLineOptionProvider merge cmdline options with defaults
type CommandLineOptionProvider struct {
	Options CommandLineOption
}

//CommandLineOption defines possible options
type CommandLineOption struct {
	Freq                   int
	ElasticsearchURL       string
	Port                   int
	Path                   string
	ElasticsearchQueryFile string
}

//ReadCmdLineOptions merge defaults with cmdline parameter
func (provider *CommandLineOptionProvider) ReadCmdLineOptions() {
	flag.StringVar(&provider.Options.ElasticsearchURL, "url", provider.Options.ElasticsearchURL, "The elastic search endpoint")
	flag.IntVar(&provider.Options.Port, "port", provider.Options.Port, "The port to listen on for HTTP requests")
	flag.StringVar(&provider.Options.Path, "path", provider.Options.Path, "The path to listen on for HTTP requests")
	flag.IntVar(&provider.Options.Freq, "freq", provider.Options.Freq, "The interval in seconds in which to query elastic search")
	flag.StringVar(&provider.Options.ElasticsearchQueryFile, "elasticqueries", provider.Options.ElasticsearchQueryFile, "The path to the queries.cfg")
	flag.Parse()
}

//PrintCmdLineOptions as logs
func (provider *CommandLineOptionProvider) PrintCmdLineOptions() {
	log.Println("Config:")
	log.Println("\tUrl:", provider.Options.ElasticsearchURL)
	log.Println("\tFreq:", provider.Options.Freq)
	log.Println("\tPort:", provider.Options.Port)
	log.Println("\tElasticsearch Configuration File:", provider.Options.ElasticsearchQueryFile)
}
