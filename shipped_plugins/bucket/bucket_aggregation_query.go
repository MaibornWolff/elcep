package main

import (
	"github.com/MaibornWolff/elcep/config"
	"github.com/olivere/elastic"
	"log"
	"time"
)

var startupTime = time.Now()

type bucketAggregationQuery struct {
	name         string
	query        string
	aggregations []string
	timeKey      string
}

func Create(query config.Query, timeKey string) *bucketAggregationQuery {
	aggregationConfig, ok := query["aggregations"];
	if !ok {
		log.Fatalf("Malformed query %v, missing 'aggregations'\n", query)
	}
	aggregationSlice, ok := aggregationConfig.([] interface{})
	if !ok {
		log.Fatalf("Malformed query %v, 'aggregations' should be of type %T\n", query, aggregationSlice)
	}
	aggregations := make([]string, len(aggregationSlice))
	for index, _field := range aggregationSlice {
		field, ok2 := _field.(string)
		if !ok2 {
			log.Fatalf("Malformed query %v, %s should be a string", query, _field)
		}
		aggregations[index] = field
	}

	return &bucketAggregationQuery{
		name:         query.Name(),
		query:        query.QueryText(),
		aggregations: aggregations,
		timeKey:      timeKey,
	}
}

func (query *bucketAggregationQuery) build(elasticClient *elastic.Client) *elastic.SearchService {
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
