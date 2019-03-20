package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"math"
)

type BucketAggregationMonitor struct {
	cache   map[string]int64
	counter *prometheus.CounterVec
	query   *BucketAggregationQuery
}

func NewAggregationMonitor(query *BucketAggregationQuery) BucketAggregationMonitor {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "elcep_logs_matched_" + query.name + "_buckets",
		Help: "Aggregates logs matching " + query.query + " to buckets",
	}, query.aggregations)
	return BucketAggregationMonitor{
		cache:   make(map[string]int64),
		counter: counter,
		query:   query,
	}
}

func (monitor *BucketAggregationMonitor) Perform(client *elastic.Client) {
	response, err := monitor.query.build(client).Do(context.Background())
	if err != nil {
		log.Printf("Elastic Query for %s failed: %v\n", monitor.query.name, err)
		return
	}
	monitor.processResponse(response)
	return
}

func (monitor *BucketAggregationMonitor) processResponse(response *elastic.SearchResult) {
	expectedAggregations := monitor.query.aggregations
	if len(expectedAggregations) > 0 {
		buckets, ok := response.Aggregations.Terms(expectedAggregations[0])
		if !ok {
			fmt.Printf("Missing terms aggregation %s in response %v\n", expectedAggregations[0], response)
		} else {
			for _, bucket := range buckets.Buckets {
				monitor.processBuckets(
					bucket,
					expectedAggregations[1:],
					withLabel(prometheus.Labels{}, expectedAggregations[0], fmt.Sprintf("%v", bucket.Key)))
			}
		}
	} else {
		counter, err := monitor.counter.GetMetricWith(prometheus.Labels{})
		if err != nil {
			fmt.Printf("Could not get counter with labels: %s\n", err)
		} else {
			inc := monitor.getInc(response.Hits.TotalHits, prometheus.Labels{})
			counter.Add(inc)
		}
	}
}

func (monitor *BucketAggregationMonitor) processBuckets(container *elastic.AggregationBucketKeyItem, expectedAggregations []string, labels prometheus.Labels) {
	if len(expectedAggregations) == 0 {
		counter, err := monitor.counter.GetMetricWith(labels)
		if err != nil {
			log.Printf("Error getting the labeled counter: %s\n", err)
		} else {
			inc := monitor.getInc(container.DocCount, labels)
			counter.Add(inc)
		}
	} else {
		buckets, _ := container.Terms(expectedAggregations[0])
		for _, bucket := range buckets.Buckets {
			monitor.processBuckets(bucket, expectedAggregations[1:], withLabel(labels, expectedAggregations[0], fmt.Sprintf("%v", bucket.Key)))
		}
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
	key := ""

	for k, v := range labels {
		key = fmt.Sprintf("%s,%s=%s", key, k, v)
	}

	lastVal, ok := monitor.cache[key]
	monitor.cache[key] = value
	if ok {
		return math.Max(0, float64(value-lastVal))
	} else {
		return float64(value)
	}
}
