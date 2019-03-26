ELCEP - Elastic Log Counter Exporter for Prometheus
===================================================

[![Go Report Card](https://goreportcard.com/badge/github.com/MaibornWolff/elcep)](https://goreportcard.com/report/github.com/MaibornWolff/elcep)
[![Release](https://img.shields.io/github/release/MaibornWolff/elcep.svg?style=flat-square)](https://github.com/MaibornWolff/elcep/releases/latest)

  - [What does it do?](#what-does-it-do)
  - [How do I use it?](#how-do-i-use-it)
  - [How do I configure it?](#how-do-i-configure-it)
  - [Which plugins can I use?](#which-plugins-can-i-use)
    - [Counter](#counter)
    - [Bucket aggregation](#bucket-aggregation)
  - [Developing custom plugins](#developing-custom-plugins)

## What does it do?
ELCEP is a small go service which provides prometheus metrics based on custom lucene queries to an elastic search instance.

## How do I use it?
Most convenient is running it as the docker image published here: https://hub.docker.com/r/maibornwolff/elcep/ eg:
```
docker run maibornwolff/elcep --url <address to elastic search instance (protocol://hostname:port)>
```

To familiarize yourself with ELCEP, check out [`examples/README.md`](examples/README.md).

## How do I configure it?
ELCEP accepts the following arguments:
```
-f, --freq=30s              The elastic search polling interval
-u, --url=http://elasticsearch:9200
                            The elastic search endpoint
-p, --port=8080             The port to listen on for HTTP requests
-c, --config=config.yml     Location of the config file
    --plugin-dir=plugins    Directory containing all the plugins
    --path="/metrics"       The resource path for the prometheus endpoint
    --timekey="@timestamp"  The timekey to use for the elasticsearch queries
-v, --version               Show application version and exit.
-h, --help                  Show help and exit.
```
These arguments can also be set via environment variables:

|  environment variable  |  argument  |  shorthand  |  default value  |
|------------------------|------------|-------------|-----------------|
| `ELCEP_POLL_FREQUENCY` | `--freq`   | `-f`        | 30s             |
| `ELCEP_ELASTIC_URL`    | `--url`    | `-u`        | http://elasticsearch:9200 |
| `ELCEP_PORT`           | `--port`   | `-p`        | 8080            |
| `ELCEP_CONFIG`         | `--config` | `-c`        | config.yml      |
| `ELCEP_PLUGIN_DIR`     | `--plugin-dir` | N/A     | plugins         |
| `ELCEP_METRICS_ENDPOINT` | `--path` | N/A         | /metrics        |
| `ELCEP_TIME_KEY`       | `--time-key` | N/A       | @timestamp      |

To configure the metrics, use the config file (`config.yml`). It has the following structure:

```yaml
plugins:
  # You can give configuration for the plugins here, if necessary.
  # Note that this section is required for each plugin you want to use, even if the plugin does not need configuration.
  counter:
    someOption: "foo"
  bucket: true

metrics:
  # logical groups
  exceptions:
    # the targeted plugin
    counter:
      # Syntax 1: `name: query` (shorthand for syntax 2)
      all: "log:exception"
      # Syntax 2: `name: configObject`
      npe:
        # query is required for all queries
        # some plugins may require more configuration for each query, e.g. for bucket aggregation
        query: "log:NullPointerException"
    # now target another plugin
    bucket:
      by_type:
        query: "log:exception"
        # you can give more options specific for that plugin
        aggregations:
          - "type"
  
  images:
    counter:
      all: "log:image"
      uploaded: "Receiving new image" 
```


### Example:

Above configuration yields to the following metrics exposed:
```bash
# HELP elcep_logs_matched_exceptions_all_total Counts number of matched logs for exceptions_all
# TYPE elcep_logs_matched_exceptions_all_total counter
elcep_logs_matched_exceptions_all_total 13
# HELP elcep_logs_matched_exceptions_npe_total Counts number of matched logs for exceptions_npe
# TYPE elcep_logs_matched_exceptions_npe_total counter
elcep_logs_matched_exceptions_npe_total 0
# HELP elcep_logs_matched_exceptions_by_type_buckets Aggregates logs matching log:exception AND bucket:true to buckets
# TYPE elcep_logs_matched_exceptions_by_type_buckets counter
elcep_logs_matched_exceptions_by_type_buckets{type="0"} 2
elcep_logs_matched_exceptions_by_type_buckets{type="1"} 2
elcep_logs_matched_exceptions_by_type_buckets{type="2"} 1
elcep_logs_matched_exceptions_by_type_buckets{type="4"} 3
elcep_logs_matched_exceptions_by_type_buckets{type="5"} 1
elcep_logs_matched_exceptions_by_type_buckets{type="7"} 1
elcep_logs_matched_exceptions_by_type_buckets{type="8"} 1
elcep_logs_matched_exceptions_by_type_buckets{type="10"} 1
elcep_logs_matched_exceptions_by_type_buckets{type="12"} 1
# HELP elcep_logs_matched_images_all_total Counts number of matched logs for images_all
# TYPE elcep_logs_matched_images_all_total counter
elcep_logs_matched_images_all_total 0
# HELP elcep_logs_matched_images_uploaded_total Counts number of matched logs for images_uploaded
# TYPE elcep_logs_matched_images_uploaded_total counter
elcep_logs_matched_images_uploaded_total 0
```

The query for elastic search and the content of the metrics depends on the used plugins

## Which plugins can I use?

Out of the box, the following plugins are provided:

### Counter

The counter plugin exposes a simple Counter metric to prometheus.
It counts the total of all matched log lines since the start of ELCEP.

#### Configuration

The plugin has no global configuration.

Each query only needs a name and a query string (which is required by default anyway).

The configured query `exceptions: "log:exception"` will match all logs that contain the string "exception" in the `log`-field.
It will count up starting from `0` at program start.

### Bucket aggregation

The bucket aggregation plugin allows to aggregate the matches by a field in the logs.
You may sub-aggregate by more fields, if necessary.
Please be aware of an exponential grow in the number of buckets when you use multiple aggregations.

#### Configuration

Each query needs an `aggregation` configured.
The configuration for a query might look like this:
```yaml
my_query:
    query: "log:searchstring"
    aggregations: ["microservice"]
```
The resulting metric will be a vector, grouping the count by the "microservice"-field of the logs.

## Developing custom plugins

Please refer to the [custom plugin guide](BUILD-CUSTOM-MONITOR.md).
