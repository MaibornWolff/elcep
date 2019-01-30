package monitor

import (
	"github.com/MaibornWolff/elcep/config"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
)

//Plugin should be used to setup prometheus
type Plugin interface {
	BuildMetrics([]config.Query) []prometheus.Collector
	Perform(*elastic.Client)
}
