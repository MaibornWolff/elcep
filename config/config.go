// This package handles all the configuration things.
package config

import (
	"io/ioutil"
	"log"
	"net/url"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Options holds all the application configuration
type Options struct {
	Freq             time.Duration
	ElasticsearchURL *url.URL
	Port             int
	Config           string
	PluginDir        string
	Path             string
	TimeKey          string
}

// Configuration holds both the application config and the plugin configuration
type Configuration struct {
	Options Options
	plugins map[string]*PluginConfig
}

// ForPlugin get the configuration for the plugin name
func (conf *Configuration) ForPlugin(name string) *PluginConfig {
	return conf.plugins[name]
}

// ReadConfig reads the configuration from the CLI options, the ENV vars and the config file
func ReadConfig() Configuration {
	options := parseCli()

	// TODO remove config file location / plugin dir
	file, err := ioutil.ReadFile(options.Config)
	if err != nil {
		log.Fatalln("Could not read config file")
	}

	plugins := parseConfigFile(file)
	return Configuration{
		Options: options,
		plugins: plugins,
	}
}

func parseCli() Options {
	freq := kingpin.Flag("freq", "The elastic search polling interval").
		Short('f').Default("30s").
		Envar("ELCEP_POLL_FREQENCY").Duration()
	elastic := kingpin.Flag("url", "The elastic search endpoint").
		Short('u').Default("http://elasticsearch:9200").
		Envar("ELCEP_ELASTIC_URL").URL()
	port := kingpin.Flag("port", "The port to listen on for HTTP requests").
		Short('p').Default("8080").
		Envar("ELCEP_PORT").Int()
	config := kingpin.Flag("config", "Location of the config file").
		Short('c').Default("config.yaml").
		Envar("ELCEP_CONFIG").ExistingFile()
	pluginDir := kingpin.Flag("plugin-dir", "Directory containing all the plugins").
		Default("plugins").
		Envar("ELCEP_PLUGIN_DIR").ExistingDir()
	path := kingpin.Flag("path", "The resource path for the prometheus endpoint").
		Default("/metrics").
		Envar("ELCEP_METRICS_ENDPOINT").String()
	timekey := kingpin.Flag("time-key", "The timekey to use for the elasticsearch queries").
		Default("@timestamp").
		Envar("ELCEP_TIME_KEY").String()

	// Todo read version during compilation: https://blog.alexellis.io/inject-build-time-vars-golang/
	kingpin.Version("0.7") // elcep version
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.Parse()

	return Options{
		Freq:             *freq,
		ElasticsearchURL: *elastic,
		Port:             *port,
		Config:           *config,
		PluginDir:        *pluginDir,
		Path:             *path,
		TimeKey:          *timekey,
	}
}
