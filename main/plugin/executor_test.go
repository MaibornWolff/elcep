package plugin

import (
	"github.com/MaibornWolff/elcep/main/config"
	"github.com/MaibornWolff/elcep/main/plugin/mock_plugin"
	"github.com/golang/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestExecutor_BuildPlugins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configuration := config.Configuration{}
	pluginConfig := config.PluginConfig{}

	mockPlugin := mock_plugin.NewMockPlugin(ctrl)
	mockPlugin.EXPECT().BuildMetrics(gomock.Eq(pluginConfig.Queries))

	executor := Executor{}
	executor.BuildPlugins(configuration, pluginConfig, func(options config.Options, pluginOptions interface{}) Plugin {
		assert.Equal(t, options, configuration.Options)
		assert.Equal(t, pluginOptions, pluginConfig.Options)
		return mockPlugin
	})
}
