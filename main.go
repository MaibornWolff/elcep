package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MaibornWolff/elcep/adapter"
	"github.com/MaibornWolff/elcep/config"
	"github.com/MaibornWolff/elcep/monitor"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	configuration := config.ReadConfig()
	executor := initExecutor(&configuration)

	go executor.PerformMonitors(configuration.Options.Freq)

	http.Handle(configuration.Options.Path, promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(configuration.Options.Port), nil))
}

func initExecutor(configuration *config.Configuration) *monitor.Executor {
	pluginProvider := adapter.NewPluginProvider("./plugins")

	elProvider := &adapter.ElasticSearchProvider{
		URL: configuration.Options.ElasticsearchURL,
	}
	client, err := elastic.NewClient(elastic.SetURL(configuration.Options.ElasticsearchURL.String()))
	if err != nil {
		log.Fatal(err)
	}
	executor := &monitor.Executor{
		// TODO remove elProvider (?)
		QueryExecution: elProvider.ExecRequest,
		ElasticClient:  client,
	}

	for name, newMon := range pluginProvider.Monitors {
		conf := configuration.ForPlugin(name)
		if conf == nil {
			log.Fatalf("Missing config for plugin %s\n", name)
		}
		executor.BuildMonitors(configuration.Options.TimeKey, *conf, newMon)
	}

	return executor
}
