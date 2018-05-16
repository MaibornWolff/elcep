package main

import (
	"time"

	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

//LogCounterMonitor for all Monitor based on Counter
type LogCounterMonitor struct {
	Query     monitor.Query
	LastCount *float64
	metrics struct {
		matchCounter         prometheus.Counter
		rpcDurationHistogram prometheus.Histogram
	}
}

//BuildMetrics should setup the Prometheus metrics
func (logMon *LogCounterMonitor) BuildMetrics(query monitor.Query) *[]prometheus.Collector {
	logMon.LastCount = new(float64)
	logMon.Query = query

	logMon.metrics.matchCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "elcep_logs_matched_" + logMon.Query.Name + "_total",
		Help: "Counts number of matched logs for " + logMon.Query.Name,
	})
	logMon.metrics.rpcDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "elcep_logs_matched_rpc_durations_" + logMon.Query.Name + "_histogram_seconds",
		Help:    "Logs matched RPC latency distributions for " + logMon.Query.Name,
		Buckets: prometheus.DefBuckets,
	})

	return &[]prometheus.Collector{logMon.metrics.matchCounter, logMon.metrics.rpcDurationHistogram}
}

//Perform executes the query for this monitor
func (logMon *LogCounterMonitor) Perform() {
	increment := logMon.countLogs()

	if increment < 0 {
		increment = 0
	}
	logMon.metrics.matchCounter.Add(increment)
}

func (logMon *LogCounterMonitor) countLogs() float64 {
	start := time.Now()
	response, _ := logMon.Query.Exec(logMon.Query.BuildBody("0", time.Now()))
	end := time.Now()

	duration := end.Sub(start).Seconds()
	logMon.metrics.rpcDurationHistogram.Observe(duration)

	increment := response.Total - *logMon.LastCount
	*logMon.LastCount = response.Total

	return increment
}
