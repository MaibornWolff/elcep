package main

import (
	"github.com/MaibornWolff/elcep/config"
	"github.com/olivere/elastic"
	"log"
)

type BucketAggregationQuery struct {
	name         string;
	query        string;
	aggregations map[string]string;
}

func Create(query config.Query) *BucketAggregationQuery {
	aggregations := make(map[string] string)

	if queryAggs, ok := query["aggregations"]; !ok {
		log.Fatalf("Malformed query %v, missing 'aggregations'\n", query)
	} else if queryAggMap, ok := queryAggs.(map[interface{}] interface{}); !ok {
		log.Fatalf("Malformed query %v, 'aggregations' should be of type %T\n", query, queryAggMap)
	} else {
		for _name, _field := range queryAggMap {
			name, ok1 := _name.(string)
			field, ok2 := _field.(string)
			if !ok1 || !ok2 {
				log.Fatalf("Malformed query %v, %s and %s should be strings", query, _name, _field)
			}
			aggregations[name] = field
		}
	}

	return &BucketAggregationQuery{
		name: query.Name(),
		query: query.QueryText(),
		aggregations: aggregations,
	}
}

func (query *BucketAggregationQuery) build(elasticClient *elastic.Client) *elastic.SearchService {
	service := elasticClient.Search().Query(elastic.NewBoolQuery())
	for aggregationName, aggregationField := range query.aggregations {
		service = service.Aggregation(aggregationName, elastic.NewTermsAggregation().Field(aggregationField))
	}
	return service
}
