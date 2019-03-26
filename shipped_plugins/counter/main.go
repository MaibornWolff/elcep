package main

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/MaibornWolff/elcep/config"
	"github.com/MaibornWolff/elcep/plugin"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
)

// LogCounterMonitor is a monitor for a certain query and exports both a Counter and a Histogram to Prometheus
type LogCounterMonitor struct {
	query     config.Query
	lastCount *int64
	metrics   struct {
		matchCounter         prometheus.Counter
		rpcDurationHistogram prometheus.Histogram
	}
}

// CounterPlugin is the exported plugin type. It implements plugin.Plugin
type CounterPlugin struct {
	monitors   []*LogCounterMonitor
	collectors []prometheus.Collector
}

func (cp *CounterPlugin) BuildMetrics(queries []config.Query) []prometheus.Collector {
	for _, query := range queries {
		log.Printf("Query: %#v\n", query)
		monitor := LogCounterMonitor{
			query:     query,
			lastCount: nil,
		}
		cp.monitors = append(cp.monitors, &monitor)
		cp.collectors = append(cp.collectors, monitor.BuildMetrics()...)
	}
	return cp.collectors
}

func (cp *CounterPlugin) Perform(elasticClient *elastic.Client) {
	for _, monitor := range cp.monitors {
		monitor.Perform(elasticClient)
	}
}

// NewPlugin must be exported. The name should be exactly "NewMonitor" and returns an instance of the custommonitor
// noinspection GoUnusedExportedFunction
func NewPlugin(_ config.Options, _ interface{}) plugin.Plugin {
	return &CounterPlugin{}
}

// BuildMetrics must exist and return a list of prometheus metrics instances
func (logMon *LogCounterMonitor) BuildMetrics() []prometheus.Collector {
	logMon.metrics.matchCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elcep_logs_matched_" + logMon.query.Name() + "_total",
		Help: "Counts number of matched logs for " + logMon.query.Name(),
	})
	logMon.metrics.rpcDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "elcep_logs_matched_rpc_durations_" + logMon.query.Name() + "_histogram_seconds",
		Help:    "Logs matched RPC latency distributions for " + logMon.query.Name(),
		Buckets: prometheus.DefBuckets,
	})

	return []prometheus.Collector{logMon.metrics.matchCounter, logMon.metrics.rpcDurationHistogram}
}

// Perform must exist and implement some custom action which runs frequently
func (logMon *LogCounterMonitor) Perform(elasticClient *elastic.Client) {
	increment, duration := logMon.runQuery(elasticClient)
	logMon.metrics.rpcDurationHistogram.Observe(duration)
	logMon.metrics.matchCounter.Add(increment)
}

func (logMon *LogCounterMonitor) runQuery(elasticClient *elastic.Client) (increment float64, duration float64) {
	start := time.Now()
	query := elastic.NewBoolQuery().Must(elastic.NewQueryStringQuery(logMon.query.QueryText()))
	response, err := elasticClient.Search().Query(query).Do(context.Background())
	duration = time.Now().Sub(start).Seconds()

	increment = 0
	if err != nil {
		log.Printf("Error on query: %#v\n", err)
		return
	}
	if logMon.lastCount == nil {
		// skip on first query
		logMon.lastCount = new(int64)
	} else {
		increment = math.Max(0, float64(response.Hits.TotalHits-*logMon.lastCount))
	}
	*logMon.lastCount = response.Hits.TotalHits
	return
}

func main() {}
