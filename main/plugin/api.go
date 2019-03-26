package plugin

import (
	"github.com/MaibornWolff/elcep/main/config"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
)

// Plugin can be implemented to extend ELCEPs functionality.
// see BUILD-CUSTOM-MONITOR.md for details
type Plugin interface {
	BuildMetrics([]config.Query) []prometheus.Collector
	Perform(*elastic.Client)
}
