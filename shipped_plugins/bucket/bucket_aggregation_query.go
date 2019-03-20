package main

import (
	"github.com/MaibornWolff/elcep/config"
	"github.com/olivere/elastic"
	"log"
	"time"
)

var startupTime = time.Now()

type BucketAggregationQuery struct {
	name         string
	query        string
	aggregations []string
	timeKey      string
}

func Create(query config.Query, timeKey string) *BucketAggregationQuery {
	aggregations := make([]string, 0, 4)

	if aggregationConfig, ok := query["aggregations"]; !ok {
		log.Fatalf("Malformed query %v, missing 'aggregations'\n", query)
	} else if aggregationSlice, ok := aggregationConfig.([] interface{}); !ok {
		log.Fatalf("Malformed query %v, 'aggregations' should be of type %T\n", query, aggregationSlice)
	} else {
		for _, _field := range aggregationSlice {
			field, ok2 := _field.(string)
			if !ok2 {
				log.Fatalf("Malformed query %v, %s should be a string", query, _field)
			}
			aggregations = append(aggregations, field)
		}
	}

	return &BucketAggregationQuery{
		name:         query.Name(),
		query:        query.QueryText(),
		aggregations: aggregations,
		timeKey:      timeKey,
	}
}

func (query *BucketAggregationQuery) build(elasticClient *elastic.Client) *elastic.SearchService {
	service := elasticClient.Search().Query(elastic.NewBoolQuery().
		Must(elastic.
			NewQueryStringQuery(query.query)).
		Filter(elastic.
			NewRangeQuery(query.timeKey).
			Gte(startupTime.Format("2006-01-02 15:04:05")).
			Format("yyyy-MM-dd HH:mm:ss")))
	if len(query.aggregations) > 0 {
		service = service.Aggregation(createAggregations(query.aggregations))
	}
	return service
}

func createAggregations(aggregationKeys []string) (string, elastic.Aggregation) {
	switch len(aggregationKeys) {
	case 0:
		log.Panicf("Cannot create aggregation without aggregation keys")
		return "", nil
	case 1:
		return aggregationKeys[0], elastic.NewTermsAggregation().Field(aggregationKeys[0])
	default:
		return aggregationKeys[0], elastic.NewTermsAggregation().Field(aggregationKeys[0]).
			SubAggregation(createAggregations(aggregationKeys[1:]))
	}
}
