ELCEP - Elastic Log Counter Exporter for Prometheus
===================================================

## What does it do?
It is a small go service which provides prometheus counter metrics based on custom lucene queries to an elastic search instance.

## How do I use it?
Configure the queries one per line in the queries.cfg in the following notation: `<name>=<query>`

Via command line arguments the following options can be overwritten:
```
  -elasticqueries string
    	The path to the queries.cfg (default "./conf/queries.cfg")
  -freq int
    	The interval in seconds in which to query elastic search (default 30)
  -host string
    	The elastic search endpoint (default "http://elasticsearch:9200")
  -path string
    	The path to listen on for HTTP requests (default "/metrics")
  -port int
    	The port to listen on for HTTP requests (default 8080)
```

### Example:
Providing this line in queries.cfg: 
```
all_application_exceptions=message:exception AND service_name:application_*
```

Will result in exposing the following metric:
```
# HELP logs_matched_all_application_exceptions_total Counts number of matched logs for all_application_exceptions
# TYPE logs_matched_all_application_exceptions_total counter
logs_matched_all_application_exceptions_total 0
```

Using that elastic search query:
```
GET /_search
{
   "query": {
       "query_string": {
           "query": "message:exception AND service_name:application_*"
       }
   },
   "size":0
}
```