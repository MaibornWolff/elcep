package plugin

import (
	"log"
	"reflect"
	"time"

	"github.com/olivere/elastic"

	"github.com/MaibornWolff/elcep/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Executor for controlling Plugin instances
type Executor struct {
	Plugins       []Plugin
	ElasticClient *elastic.Client
}

// BuildPlugins create new Instances form given Monitortype for each query and register all metrics
func (executor *Executor) BuildPlugins(configuration config.PluginConfig, newPlugin func(interface{}) Plugin) {
	plugin := newPlugin(configuration.Options)
	metrics := plugin.BuildMetrics(configuration.Queries)
	executor.register(metrics)
	executor.Plugins = append(executor.Plugins, plugin)
	log.Println("Plugin loaded: ", reflect.TypeOf(plugin))
}

// RunPlugins runs all Plugins in a loop
func (executor *Executor) RunPlugins(freq time.Duration) {
	nextExecution := time.Now()
	for {
		time.Sleep(nextExecution.Sub(time.Now()))
		nextExecution = nextExecution.Add(freq)

		for _, plugin := range executor.Plugins {
			plugin.Perform(executor.ElasticClient)
		}
	}
}

func (executor *Executor) register(metrics []prometheus.Collector) {
	for _, metric := range metrics {
		prometheus.MustRegister(metric)
	}
}
