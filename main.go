package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MaibornWolff/elcep/adapter"
	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cmdlineProvider := &adapter.CommandLineOptionProvider{
		Options: adapter.CommandLineOption{
			Freq:             30,
			ElasticsearchURL: "http://elasticsearch:9200",
			Port:             8080,
			Path:             "/metrics",
			ElasticsearchQueryFile: "./conf/queries.cfg",
		},
	}
	cmdlineProvider.ReadCmdLineOptions()
	cmdlineProvider.PrintCmdLineOptions()

	elProvider := &adapter.ElasticSearchProvider{
		URL: cmdlineProvider.Options.ElasticsearchURL,
	}

	executor := &monitor.Executor{}
	queryProvider := &adapter.QueryProvider{}

	executor.QueryExecution = elProvider.ExecRequest
	queryProvider.Read(cmdlineProvider.Options.ElasticsearchQueryFile)
	queryProvider.Print()

	//TODO if we want to use muliple monitor types here we have to distict the queries
	executor.BuildMonitors(queryProvider.Queries, func() monitor.LogMonitor {
		return &LogCounterMonitor{}
	})

	//Example: Register some other monitor (should get other queries)
	//executor.BuildMonitors(queryProvider.Queries, func() monitor.LogMonitor {
	//	return &SomeOtherMonitor{}
	//})

	go executor.PerformMonitors(cmdlineProvider.Options.Freq)

	http.Handle(cmdlineProvider.Options.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cmdlineProvider.Options.Port), nil))
}
