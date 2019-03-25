package main

import (
	"context"
	"fmt"
	"github.com/mitchellh/hashstructure"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math"
)

type bucketAggregationMonitor struct {
	cache   map[uint64]int64
	counter *prometheus.CounterVec
	query   *bucketAggregationQuery
}

func NewAggregationMonitor(query *bucketAggregationQuery) *bucketAggregationMonitor {
	return &bucketAggregationMonitor{
		cache: make(map[uint64]int64),
		counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "elcep_logs_matched_" + query.name + "_buckets",
			Help: "Aggregates logs matching " + query.query + " to buckets",
		}, query.aggregations),
		query: query,
	}
}

func (monitor *bucketAggregationMonitor) Perform(client *elastic.Client) {
	response, err := monitor.query.build(client).Do(context.Background())
	if err != nil {
		log.Printf("Elastic Query for %s failed: %v\n", monitor.query.name, err)
		return
	}
	monitor.processAggregations(
		&response.Aggregations, monitor.query.aggregations, prometheus.Labels{}, response.Hits.TotalHits)
}

func (monitor *bucketAggregationMonitor) processAggregations(
	container *elastic.Aggregations, expectedAggregations []string, labels prometheus.Labels, hits int64) {

	if len(expectedAggregations) == 0 {
		monitor.updateCounter(hits, labels)
		return
	}
	if buckets, ok := container.Terms(expectedAggregations[0]); !ok {
		log.Printf("Missing terms aggregation %s in response %v\n", expectedAggregations[0], container)
	} else {
		for _, bucket := range buckets.Buckets {
			monitor.processAggregations(
				&bucket.Aggregations,
				expectedAggregations[1:],
				withLabel(labels, expectedAggregations[0], fmt.Sprintf("%v", bucket.Key)),
				bucket.DocCount)
		}
	}
}

func (monitor *bucketAggregationMonitor) updateCounter(newCount int64, labels prometheus.Labels) {
	if counter, err := monitor.counter.GetMetricWith(labels); err != nil {
		log.Printf("Error getting the labeled counter: %s\n", err)
	} else {
		counter.Add(monitor.getInc(newCount, labels))
	}
}

func withLabel(labels prometheus.Labels, key string, value string) prometheus.Labels {
	newLabels := prometheus.Labels{
		key: value,
	}
	for k, v := range labels {
		newLabels[k] = v
	}
	return newLabels
}

func (monitor *bucketAggregationMonitor) getInc(value int64, labels prometheus.Labels) float64 {
	hash, err := hashstructure.Hash(labels, nil)
	if err != nil {
		log.Printf("An error occured while hashing the labels: %v\n", err)
	}

	lastVal, ok := monitor.cache[hash]
	monitor.cache[hash] = value
	if ok {
		return math.Max(0, float64(value-lastVal))
	} else {
		return float64(value)
	}
}
