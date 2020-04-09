package main

import (
	"github.com/MaibornWolff/elcep/main/config"
	"github.com/olivere/elastic"
	"github.com/patrickmn/go-cache"
	"log"
	"regexp"
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

		aggregations[index] = getAllowedPrometheusLabel(field)
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
			Format("yyyy-MM-dd HH:mm:ss"))).
		FilterPath("hits.total,aggregations")
	if len(query.aggregations) > 0 {
		originAggregations := getOriginalAggregationKeys(query.aggregations)
		service = service.Aggregation(createAggregations(originAggregations))
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

// Retrieve the origin label name from the cache
func getOriginalAggregationKeys(s []string) []string{
	strSlice := make([]string, 0)
	for _ , str := range s {
		if x, found := bucketCache.Get(str); found {
			value := x.(string)
			strSlice = append(strSlice, value)
		}
	}
	return strSlice
}


// Replace given string special characters and return '_' instead.
// save the origin label to cache.
// https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels
func getAllowedPrometheusLabel(s string) string {
	var re = regexp.MustCompile(`[!@#$%^&*(),./\\?":{}|<>]`)
	r := re.ReplaceAllString(s, `${1}_${2}`)
	bucketCache.Set(r, s, cache.NoExpiration) // set origin label to the cache

	return  r
}