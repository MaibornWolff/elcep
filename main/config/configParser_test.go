package config

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func Test_Read_Queries_From_Files(t *testing.T) {
	config := parseConfigFile([]byte(configFile))

	assert.NotEqual(t, config["counter"], nil)

	var counterConf *PluginConfig
	counterConf = config["counter"]
	assert.NotEqual(t, counterConf.Options, nil)
	assert.Equal(t, counterConf.Options.(map[interface{}] interface{})["someOption"], "bla")

	assert.Equal(t, len(counterConf.Queries), 4)
}

const configFile = `---
plugins:
  # you can give configuration for the plugins here, if necessary
  counter:
    someOption: "bla"

metrics:
  exceptions:
    counter:
      all: "log:exception"
      npe:
        query: "log:NullPointerException"
  
  images:
    counter:
      all: "log:image"
      uploaded: "Receiving new image"
`
