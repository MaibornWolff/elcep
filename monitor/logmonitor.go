package monitor

import "github.com/prometheus/client_golang/prometheus"

//LogMonitor should be used to setup prometheus
type LogMonitor interface {
	BuildMetrics(Query) *[]prometheus.Collector
	Perform()
}
