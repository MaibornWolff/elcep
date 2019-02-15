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

// BuildPlugins creates the plugin instances
func (executor *Executor) BuildPlugins(configuration config.Configuration, pluginConfig config.PluginConfig, newPlugin factoryMethodType) {
	plugin := newPlugin(configuration.Options, pluginConfig.Options)
	metrics := plugin.BuildMetrics(pluginConfig.Queries)
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
