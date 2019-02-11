package main

import (
	"context"
	"log"
	"time"

	"github.com/olivere/elastic"

	"github.com/MaibornWolff/elcep/config"
	"github.com/MaibornWolff/elcep/plugin"
	"github.com/prometheus/client_golang/prometheus"
)

var startupTime = time.Now()

// LogCounterMonitor is in this example the exported plugin type. It must implement BuildMetrics and Perform like below
type LogCounterMonitor struct {
	gauge     prometheus.Gauge
	query     config.Query
	LastCount *float64
	metrics   struct {
		matchCounter         prometheus.Counter
		rpcDurationHistogram prometheus.Histogram
	}
}

type CounterPlugin struct {
	monitors   []*LogCounterMonitor
	collectors []prometheus.Collector
}

func (cp *CounterPlugin) BuildMetrics(queries []config.Query) []prometheus.Collector {
	for _, query := range queries {
		log.Printf("Query: %#v\n", query)
		log.Printf(" - QText: %#v\n\n", query.QueryText())
		monitor := LogCounterMonitor{}
		cp.monitors = append(cp.monitors, &monitor)
		cp.collectors = append(cp.collectors, monitor.BuildMetrics(query)...)
	}
	return cp.collectors
}

func (cp *CounterPlugin) Perform(elasticClient *elastic.Client) {
	for _, monitor := range cp.monitors {
		monitor.Perform(elasticClient)
	}
}

// NewPlugin must be exported. The name should be exactly "NewMonitor" and returns an instance of the custommonitor
func NewPlugin(config interface{}) plugin.Plugin {
	return &CounterPlugin{}
}

// BuildMetrics must exist and return a list of prometheus metrics instances
func (logMon *LogCounterMonitor) BuildMetrics(query config.Query) []prometheus.Collector {
	logMon.LastCount = new(float64)
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
func (logMon *LogCounterMonitor) Perform(elasticClient *elastic.Client) {
	increment, duration := logMon.runQuery(elasticClient)
	logMon.metrics.rpcDurationHistogram.Observe(duration)
	logMon.metrics.matchCounter.Add(increment)
}

func (logMon *LogCounterMonitor) runQuery(elasticClient *elastic.Client) (increment float64, duration float64) {
	start := time.Now()
	query := elastic.NewBoolQuery().
		Must(elastic.
			NewQueryStringQuery(logMon.query.QueryText())).
		Filter(elastic.
			NewRangeQuery("@timestamp").
			Gte(startupTime.Format("2006-01-02 15:04:05")).
			Format("yyyy-MM-dd HH:mm:ss"))
	response, err := elasticClient.Search().Query(query).Do(context.Background())
	duration = time.Now().Sub(start).Seconds()

	if err == nil {
		increment = float64(response.Hits.TotalHits) - *logMon.LastCount
		*logMon.LastCount += increment
	} else {
		log.Printf("Error on query: %#v\n", err)
		increment = 0
	}
	return
}

func main() {}
