package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MaibornWolff/elcep/adapter"
	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

var startup_time = time.Now()

func main() {
	cmdlineProvider := &adapter.CommandLineOptionProvider{
		Options: adapter.CommandLineOption{
			Freq:                   30,
			ElasticsearchURL:       "http://elasticsearch:9200",
			Port:                   8080,
			Path:                   "/metrics",
			ElasticsearchQueryFile: "./conf/queries.cfg",
			TimeKey:                "@timestamp",
		},
	}
	cmdlineProvider.ReadCmdLineOptions()
	cmdlineProvider.PrintCmdLineOptions()

	elProvider := &adapter.ElasticSearchProvider{
		URL: cmdlineProvider.Options.ElasticsearchURL,
	}

	pluginProvider := adapter.NewPluginProvider("./plugins")

	executor := &monitor.Executor{}
	queryProvider := adapter.NewQueryProvider(cmdlineProvider.Options.ElasticsearchQueryFile, pluginProvider.GetPluginNames())

	executor.QueryExecution = elProvider.ExecRequest
	queryProvider.Print()

	buildAllMonitors(executor, pluginProvider, queryProvider, cmdlineProvider)

	go executor.PerformMonitors(cmdlineProvider.Options.Freq)

	http.Handle(cmdlineProvider.Options.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(cmdlineProvider.Options.Port), nil))
}

func buildAllMonitors(executor *monitor.Executor, pluginProvider *adapter.PluginProvider, queryProvider *adapter.QueryProvider, cmdlineProvider *adapter.CommandLineOptionProvider) {
	executor.BuildMonitors(cmdlineProvider.Options.TimeKey, queryProvider.QuerySets["default"].Queries, func() monitor.LogMonitor {
		return &LogCounterMonitor{}
	})

	for name, newMon := range pluginProvider.Monitors {
		executor.BuildMonitors(cmdlineProvider.Options.TimeKey, queryProvider.QuerySets[name].Queries, newMon)
	}
}
