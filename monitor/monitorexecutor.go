package monitor

import (
	"log"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//Executor for controlling LogMonitor instances
type Executor struct {
	LogMonitors    []LogMonitor
	QueryExecution func(string) (*Hits, error)
}

//BuildMonitors create new Instances form given Monitortype for each query and register all metrics
func (executor *Executor) BuildMonitors(plainqueries map[string]string, newMonitor func() LogMonitor) {
	for name, queryBody := range plainqueries {
		logMonitor := newMonitor()
		query := Query{
			Name:   name,
			filter: queryBody,
			Exec:   executor.QueryExecution,
		}

		metrics := logMonitor.BuildMetrics(query)
		executor.register(metrics)
		executor.LogMonitors = append(executor.LogMonitors, logMonitor)

		log.Println("Register Monitor:", reflect.TypeOf(logMonitor), " for query:", query.Name)
	}
}

//PerformMonitors runs all Monitors in a loop
func (executor *Executor) PerformMonitors(freq int) {
	for {
		for _, logMon := range executor.LogMonitors {
			logMon.Perform()
		}
		time.Sleep(time.Duration(freq) * time.Second)
	}
}

func (executor *Executor) register(metrics *[]prometheus.Collector) {
	for _, metric := range *metrics {
		prometheus.MustRegister(metric)
	}
}
