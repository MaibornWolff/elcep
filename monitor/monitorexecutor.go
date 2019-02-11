package monitor

import (
	"log"
	"reflect"
	"time"

	"github.com/olivere/elastic"

	"github.com/MaibornWolff/elcep/config"
	"github.com/prometheus/client_golang/prometheus"
)

//Executor for controlling LogMonitor instances
type Executor struct {
	Plugins        []Plugin
	ElasticClient  *elastic.Client
}

//BuildMonitors create new Instances form given Monitortype for each query and register all metrics
func (executor *Executor) BuildMonitors(timeKey string, configuration config.PluginConfig, newMonitor func(interface{}) Plugin) {
	plugin := newMonitor(configuration.Options)
	metrics := plugin.BuildMetrics(configuration.Queries)
	executor.register(metrics)
	executor.Plugins = append(executor.Plugins, plugin)
	log.Println("Plugin loaded:", reflect.TypeOf(plugin))
}

//PerformMonitors runs all Monitors in a loop
func (executor *Executor) PerformMonitors(freq time.Duration) {
	for {
		for _, plugin := range executor.Plugins {
			plugin.Perform(executor.ElasticClient)
		}
		// TODO subtract passed time or put plugin.Perform in goroutine
		time.Sleep(freq)
	}
}

func (executor *Executor) register(metrics []prometheus.Collector) {
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}
