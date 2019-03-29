package main

import (
	"github.com/MaibornWolff/elcep/main/config"
	"github.com/MaibornWolff/elcep/main/plugin"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

// The factory method for the plugin
// noinspection GoUnusedExportedFunction
func NewPlugin(options config.Options, _ interface{}) plugin.Plugin {
	return &bucketAggregationPlugin{
		timeKey: options.TimeKey,
	}
}

type bucketAggregationPlugin struct {
	timeKey    string
	monitors   []*bucketAggregationMonitor
	collectors []prometheus.Collector
}

func (plugin *bucketAggregationPlugin) BuildMetrics(queries []config.Query) []prometheus.Collector {
	for _, query := range queries {
		log.Printf("Query loaded: %#v\n", query)
		monitor := NewAggregationMonitor(Create(query, plugin.timeKey))
		plugin.monitors = append(plugin.monitors, monitor)
		plugin.collectors = append(plugin.collectors, monitor.counter)
	}
	return plugin.collectors
}

func (plugin *bucketAggregationPlugin) Perform(elasticClient *elastic.Client) {
	for _, monitor := range plugin.monitors {
		monitor.Perform(elasticClient)
	}
}

func main() {}
