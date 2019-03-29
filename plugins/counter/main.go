package main

import (
	"context"
	"log"
	"time"

	"github.com/MaibornWolff/elcep/main/config"
	"github.com/MaibornWolff/elcep/main/plugin"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
)

var startupTime = time.Now()

// LogCounterMonitor is a monitor for a certain query and exports both a Counter and a Histogram to Prometheus
type LogCounterMonitor struct {
	query     config.Query
	LastCount *int64
	metrics   struct {
		matchCounter         prometheus.Counter
		rpcDurationHistogram prometheus.Histogram
	}
}

// CounterPlugin is the exported plugin type. It implements plugin.Plugin
type CounterPlugin struct {
	timeKey    string
	monitors   []*LogCounterMonitor
	collectors []prometheus.Collector
}

func (cp *CounterPlugin) BuildMetrics(queries []config.Query) []prometheus.Collector {
	for _, query := range queries {
		log.Printf("Query loaded: %#v\n", query)
		monitor := LogCounterMonitor{}
		cp.monitors = append(cp.monitors, &monitor)
		cp.collectors = append(cp.collectors, monitor.BuildMetrics(query)...)
	}
	return cp.collectors
}

func (cp *CounterPlugin) Perform(elasticClient *elastic.Client) {
	for _, monitor := range cp.monitors {
		monitor.Perform(elasticClient, cp.timeKey)
	}
}

// NewPlugin must be exported. The name should be exactly "NewMonitor" and returns an instance of the custommonitor
// noinspection GoUnusedExportedFunction
func NewPlugin(options config.Options, _ interface{}) plugin.Plugin {
	return &CounterPlugin{
		timeKey: options.TimeKey,
	}
}

// BuildMetrics must exist and return a list of prometheus metrics instances
func (logMon *LogCounterMonitor) BuildMetrics(query config.Query) []prometheus.Collector {
	logMon.LastCount = new(int64)
	logMon.query = query

	logMon.metrics.matchCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elcep_logs_matched_" + query.Name() + "_total",
		Help: "Counts number of matched logs for " + query.Name(),
	})
	logMon.metrics.rpcDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "elcep_logs_matched_rpc_durations_" + query.Name() + "_histogram_seconds",
		Help:    "Logs matched RPC latency distributions for " + query.Name(),
		Buckets: prometheus.DefBuckets,
	})

	return []prometheus.Collector{logMon.metrics.matchCounter, logMon.metrics.rpcDurationHistogram}
}

// Perform must exist and implement some custom action which runs frequently
func (logMon *LogCounterMonitor) Perform(elasticClient *elastic.Client, timeKey string) {
	increment, duration := logMon.runQuery(elasticClient, timeKey)
	logMon.metrics.rpcDurationHistogram.Observe(duration)
	logMon.metrics.matchCounter.Add(float64(increment))
}

func (logMon *LogCounterMonitor) runQuery(elasticClient *elastic.Client, timeKey string) (increment int64, duration float64) {
	start := time.Now()
	query := elastic.NewBoolQuery().
		Must(elastic.
			NewQueryStringQuery(logMon.query.QueryText())).
		Filter(elastic.
			NewRangeQuery(timeKey).
			Gte(startupTime.Format("2006-01-02 15:04:05")).
			Format("yyyy-MM-dd HH:mm:ss"))
	count, err := elasticClient.Count().Query(query).Do(context.Background())
	duration = time.Now().Sub(start).Seconds()

	if err == nil {
		increment = count - *logMon.LastCount
		if increment < 0 {
			increment = 0
		}
		*logMon.LastCount = count
	} else {
		log.Printf("Error on query: %#v\n", err)
		increment = 0
	}
	return
}

func main() {}
