package config

import (
	"io/ioutil"
	"log"
	"net/url"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type Options struct {
	Freq             time.Duration
	ElasticsearchURL *url.URL
	Port             int
	Config           string
	PluginDir        string
	Path             string
	TimeKey          string
}

type Configuration struct {
	Options Options
	plugins map[string]*PluginConfig
}

func (conf *Configuration) ForPlugin(name string) *PluginConfig {
	return conf.plugins[name]
}

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
	// TODO check if duration without unit works
	freq := kingpin.Flag("freq", "The elastic search polling interval").
		Short('f').Default("30s").
		OverrideDefaultFromEnvar("ELCEP_POLL_FREQENCY").Duration()
	elastic := kingpin.Flag("url", "The elastic search endpoint").
		Short('u').Default("http://elasticsearch:9200").
		OverrideDefaultFromEnvar("ELCEP_ELASTIC_URL").URL()
	port := kingpin.Flag("port", "The port to listen on for HTTP requests").
		Short('p').Default("8080").
		OverrideDefaultFromEnvar("ELCEP_PORT").Int()
	config := kingpin.Flag("config", "Location of the config file").
		Short('c').Default("config.yaml").
		OverrideDefaultFromEnvar("ELCEP_CONFIG").ExistingFile()
	pluginDir := kingpin.Flag("plugin-dir", "Directory containing all the plugins").
		Default("plugins").
		OverrideDefaultFromEnvar("ELCEP_PLUGIN_DIR").ExistingDir()
	path := kingpin.Flag("path", "The port to listen on for HTTP requests").
		Default("/metrics").
		OverrideDefaultFromEnvar("ELCEP_METRICS_ENDPOINT").String()
	timekey := kingpin.Flag("timekey", "The port to listen on for HTTP requests").
		Default("@timestamp").
		OverrideDefaultFromEnvar("ELCEP_TIME_KEY").String()

	kingpin.CommandLine.HelpFlag.Short('h')
	// TODO kingpin.CommandLine.VersionFlag.Short('v')
	// Todo read during compilation (?)
	kingpin.Version("0.7") // elcep version
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
