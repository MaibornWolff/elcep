package main

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type BucketAggregationMonitor struct {
	cache   map[string]map[string]float64
	counter *prometheus.CounterVec
	query   *BucketAggregationQuery
}

func NewAggregationMonitor(query *BucketAggregationQuery) BucketAggregationMonitor {var labels []string
	for name, _ := range query.aggregations {
		labels = append(labels, name)
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "elcep_logs_matched_" + query.name + "_buckets",
		Help: "Aggregates logs matching " + query.query + " to buckets",
	}, labels)
	return BucketAggregationMonitor{
		cache: make(map[string]map[string]float64),
		counter: counter,
		query: query,
	}
}

func (monitor *BucketAggregationMonitor) Perform(client *elastic.Client) {
	response, err := monitor.query.build(client).Do(context.Background())
	if err != nil {
		log.Printf("Elastic Query for %s failed: %v\n", monitor.query.name, err)
		return
	}

	for name := range monitor.query.aggregations {
		items, ok := response.Aggregations.Terms(name)
		if !ok {
			log.Printf("Answer missing aggregation %s\n", name)
		}
		if _, ok := monitor.cache[name]; !ok {
			monitor.cache[name] = make(map[string]float64)
		}
		for _, bucket := range items.Buckets {
			log.Printf("%#v\n", monitor.counter)
			counter, err := monitor.counter.GetMetricWith(prometheus.Labels{
				name: bucket.Key.(string),
			})
			if err != nil {
				log.Printf("Could not get metric with label '%s' '%s': %v\n", name, bucket.Key.(string), err)
			}

			lastValue, ok := monitor.cache[name][bucket.Key.(string)]
			if !ok {
				lastValue = 0
			}
			monitor.cache[name][bucket.Key.(string)] = float64(bucket.DocCount)
			inc := float64(bucket.DocCount) - lastValue
			counter.Add(inc)
		}
	}
}
