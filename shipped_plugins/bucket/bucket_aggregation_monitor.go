package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math"
	"strings"
)

type BucketAggregationMonitor struct {
	cache   map[string]int64
	counter *prometheus.CounterVec
	query   *BucketAggregationQuery
}

func NewAggregationMonitor(query *BucketAggregationQuery) BucketAggregationMonitor {
	return BucketAggregationMonitor{
		cache: make(map[string]int64),
		counter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "elcep_logs_matched_" + query.name + "_buckets",
			Help: "Aggregates logs matching " + query.query + " to buckets",
		}, query.aggregations),
		query: query,
	}
}

func (monitor *BucketAggregationMonitor) Perform(client *elastic.Client) {
	if response, err := monitor.query.build(client).Do(context.Background()); err != nil {
		log.Printf("Elastic Query for %s failed: %v\n", monitor.query.name, err)
	} else {
		monitor.processAggregations(
			&response.Aggregations, monitor.query.aggregations, prometheus.Labels{}, response.Hits.TotalHits)
	}
}

func (monitor *BucketAggregationMonitor) processAggregations(
	container *elastic.Aggregations, expectedAggregations []string, labels prometheus.Labels, hits int64) {
	if len(expectedAggregations) > 0 {
		if buckets, ok := container.Terms(expectedAggregations[0]); !ok {
			fmt.Printf("Missing terms aggregation %s in response %v\n", expectedAggregations[0], container)
		} else {
			for _, bucket := range buckets.Buckets {
				monitor.processAggregations(
					&bucket.Aggregations,
					expectedAggregations[1:],
					withLabel(labels, expectedAggregations[0], fmt.Sprintf("%v", bucket.Key)),
					bucket.DocCount)
			}
		}
	} else {
		monitor.updateCounter(hits, labels)
	}
}

func (monitor *BucketAggregationMonitor) updateCounter(newCount int64, labels prometheus.Labels) {
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

func (monitor *BucketAggregationMonitor) getInc(value int64, labels prometheus.Labels) float64 {
	keys := make([]string, len(labels))
	for i, bucket := range monitor.query.aggregations {
		keys[i] = labels[bucket]
	}
	key := strings.Join(keys, "|")

	lastVal, ok := monitor.cache[key]
	monitor.cache[key] = value
	if ok {
		return math.Max(0, float64(value-lastVal))
	} else {
		return float64(value)
	}
}
