package main

import (
	"time"

	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

var startupTime = time.Now()

// LogCounterMonitor is in this example the exported monitor type. It must implement BuildMetrics and Perform like below
type LogCounterMonitor struct {
	gauge     prometheus.Gauge
	query     monitor.Query
	LastCount *float64
	metrics   struct {
		matchCounter         prometheus.Counter
		rpcDurationHistogram prometheus.Histogram
	}
}

// NewMonitor must be exported. The name should be exactly "NewMonitor" and returns an instance of the custommonitor
func NewMonitor() monitor.LogMonitor {
	return &LogCounterMonitor{}
}

// BuildMetrics must exist and return a list of prometheus metrics instances
func (logMon *LogCounterMonitor) BuildMetrics(query monitor.Query) []prometheus.Collector {
	logMon.LastCount = new(float64)
	logMon.query = query

	logMon.metrics.matchCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elcep_logs_matched_" + logMon.query.Name + "_total",
		Help: "Counts number of matched logs for " + logMon.query.Name,
	})
	logMon.metrics.rpcDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "elcep_logs_matched_rpc_durations_" + logMon.query.Name + "_histogram_seconds",
		Help:    "Logs matched RPC latency distributions for " + logMon.query.Name,
		Buckets: prometheus.DefBuckets,
	})

	return []prometheus.Collector{logMon.metrics.matchCounter, logMon.metrics.rpcDurationHistogram}
}

// Perform must exist and implement some custom action which runs frequently
func (logMon *LogCounterMonitor) Perform() {
	increment, duration := logMon.runQuery()
	logMon.metrics.rpcDurationHistogram.Observe(duration)
	logMon.metrics.matchCounter.Add(increment)
}

func (logMon *LogCounterMonitor) runQuery() (increment float64, duration float64) {
	start := time.Now()
	response, err := logMon.query.Exec(logMon.query.BuildBody("0", startupTime))
	duration = time.Now().Sub(start).Seconds()

	if err == nil {
		increment = response.Total - *logMon.LastCount
		*logMon.LastCount = response.Total
	} else {
		increment = 0
	}
	return
}

func main() {}
