package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/MaibornWolff/elcep/adapter"
	"github.com/MaibornWolff/elcep/config"
	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	options := config.ParseCliOptions()
	options.PrintCliOptions()
	executor := initExecutor(options)

	go executor.PerformMonitors(options.Freq)

	http.Handle(options.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(options.Port), nil))
}

func initExecutor(cliOptions *config.CommandLineOption) *monitor.Executor {
	pluginProvider := adapter.NewPluginProvider("./plugins")

	elProvider := &adapter.ElasticSearchProvider{
		URL: cliOptions.ElasticsearchURL,
	}
	executor := &monitor.Executor{
		QueryExecution: elProvider.ExecRequest,
	}

	pluginConfig := config.ReadConfig(pluginProvider.GetPluginNames(), getConfigFile)
	pluginConfig.Print()

	for name, newMon := range pluginProvider.Monitors {
		executor.BuildMonitors(cliOptions.TimeKey, pluginConfig.ForPlugin(name), newMon)
	}

	return executor
}

func getConfigFile(pluginName string) io.ReadCloser {
	filepath := filepath.Join("conf", pluginName+".cfg")
	fileHandle, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Could not read config file for plugin %s, expected file: %s", pluginName, filepath)
	}
	return fileHandle
}
