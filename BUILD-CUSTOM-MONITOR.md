# How to build a Custom Monitor as Plugin

In most cases the default counter monitor is sufficient. But if you need some custom logic depending on the query data, or another Prometheus metric like Histogram, Gauge or Summary, it is possible to inject a monitor plugin in a few simple steps.

## Setup a Montitor plugin

### Step 1

Create a new file for queries inside the elcep conf folder (The name of file must be in the case `custommonitor.cfg`). Define some queries inside (see [README.md](README.md)). 

### Step 2

Create a new folder for a empty Go project under your gopath
```
$ mkdir custommonitor
$ cd custommonitor
```

### Step 3

Basically you need a new type of struct like `CustomMonitor` with 2 implemented Methods `BuildMetrics` and `Perform`. The BuildMetrics should  return one or more new Prometheus metrics. For each query in the corresponding queryfile, all metrics will be registered. As a parameter you get a query instance, which should be saved in your struct.

The Perform Method will executed in a loop depending of the frequency. Here you can execute the query with `yourmonitor.query.Exec(yourmonitor.query.BuildBody("10000"))`. The parameter is the maximum fetched rows from the elastic db. If you dont need the data rows itself, you can set the value to "0". 
Depending on your metric, to your stuff here.

The last importent thing is the `NewMonitor` func. It must return a new Instance of your custom monitor.

Full example:

```
package main

import (
	"github.com/MaibornWolff/elcep/monitor"
	"github.com/prometheus/client_golang/prometheus"
)

//CustomMonitor is in this example the exported monitor type. It must implement BuildMetrics and Perform like below
type CustomMonitor struct {
	gauge prometheus.Gauge
	query monitor.Query
}

//BuildMetrics must exist and return a list of prometheus metrics instances
func (mon *CustomMonitor) BuildMetrics(query monitor.Query) *[]prometheus.Collector {
	mon.query = query
	mon.gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "elcep_custommonitor_" + mon.query.Name + "_gauge",
		Help: "average of imported data",
	})

	return &[]prometheus.Collector{mon.gauge}
}

//Perform must exist and implement some custom action which runs frequently
func (mon *CustomMonitor) Perform() {
	response, _ := mon.query.Exec(mon.query.BuildBody("10000"))
	mon.gauge.Set(calculateAverage(response))
}

//NewMonitor must be exported. The name should be exactly "NewMonitor" and returns an instance of the custommonitor
func NewMonitor() monitor.LogMonitor {
	return &CustomMonitor{}
}

func calculateAverage(hits *monitor.Hits) float64 {
	amount := 0.0
	for _, result := range hits.Results {
		dataImported, _ := result.(map[string]interface{})["_source"].(map[string]interface{})["dataImported"].(float64)
		amount += dataImported
	}

	return amount / hits.Total
}

func main() {}
```

### Step 4

Build your plugin and copy the output inside the elcep plugins folder. (The easiest way is to use docker bind volumes for the "plugins" and "conf" folder)

```
go build --buildmode=plugin -o custommonitor.so
```

**Caution:** the plugin buildmode is currently only available on linux. If you work on mac or windows, you can use a container to build your plugin: 

```
$ docker run -v $PWD:/go/src/app  -it maibornwolff/elcep:builder-1.10.2 sh
# apk add --no-cache git gcc musl-dev
# cd /go/src/app
# go get
# go build --buildmode=plugin -o custommonitor.so
```

### step 5

Start the elcep container. Thats it!