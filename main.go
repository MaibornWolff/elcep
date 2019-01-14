package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MaibornWolff/elcep/adapter"
	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startupTime = time.Now()

func main() {
	options := adapter.ParseOptions()
	options.PrintCmdLineOptions()
	executor := initExecutor(options)

	go executor.PerformMonitors(options.Freq)

	http.Handle(options.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(options.Port), nil))
}

func initExecutor(cliOptions *adapter.CommandLineOption) *monitor.Executor {
	pluginProvider := adapter.NewPluginProvider("./plugins")

	elProvider := &adapter.ElasticSearchProvider{
		URL: cliOptions.ElasticsearchURL,
	}
	executor := &monitor.Executor{}
	executor.QueryExecution = elProvider.ExecRequest

	queryProvider := adapter.NewQueryProvider(pluginProvider.GetPluginNames())
	queryProvider.Print()

	buildAllMonitors(executor, pluginProvider, queryProvider, cliOptions)

	return executor
}

func buildAllMonitors(executor *monitor.Executor, pluginProvider *adapter.PluginProvider, queryProvider *adapter.QueryProvider, options *adapter.CommandLineOption) {
	for name, newMon := range pluginProvider.Monitors {
		executor.BuildMonitors(options.TimeKey, queryProvider.QuerySets[name].Queries, newMon)
	}
}
