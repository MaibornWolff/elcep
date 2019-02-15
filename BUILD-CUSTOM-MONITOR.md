# How to build a Custom Monitor as Plugin

In most cases the provided metrics are sufficient.
But if you need some custom logic depending on the query data, or another Prometheus metric like Histogram, Gauge or Summary, it is possible to inject a monitor plugin in a few simple steps.

## Setup a Montitor plugin

Create a new folder for a empty Go project under your gopath
```
$ mkdir custommonitor
$ cd custommonitor
```

Provide a factory method called `NewPlugin` of type `func(config.Options, interface{}) plugin.Plugin`.
The general options are passed as the first, the plugin options from the config file as the second parameter.
The function needs to return an instance of `plugin.Plugin`, that requires two methods:

```
type Plugin interface {
   	BuildMetrics([]config.Query) []prometheus.Collector
   	Perform(*elastic.Client)
}
```

First, `BuildMetrics` will be called with the queries that are configured in the config file.
Simply return all prometheus metric objects that should be exposed.
`Perform` requests you to execute the query.
You can easily do this using the [elastic search client implementation](https://olivere.github.io/elastic/) in go.
In case you struggle with the API, you can easily debug the query by enabling the wireshark container (see the provided example `docker-compose`).

For a full example have a look at the implementations at `shipped_plugins`, but the general structure should look like this:
```
package main

import(...)

type CustomPlugin struct {
    // probably you want to save some things here
}

func (plugin *CustomPlugin) BuildMetrics(queries []config.Query) *[]prometheus.Collector {
    // add some prometheus metrics here
	return &[]prometheus.Collector{ myMetric }
}

func (plugin *CustomPlugin) Perform(client *elastic.Client) {
	response, _ := client.Search().Query(...).Do(context.Background())
	myMetric.AddSomeData(calcSomething(response))
}

//NewMonitor must be exported. The name should be exactly "NewMonitor" and returns an instance of the custommonitor
func NewPlugin(options config.Options, pluginOptions interface{}) plugin.Plugin {
	return &CustomPlugin{}
}

func main() {}
```

Build your plugin and copy the output inside the ELCEP plugins folder.
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

It's as simple as that.
Go on now and create something awesome!
