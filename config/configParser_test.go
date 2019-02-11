package config

import (
	"testing"
)

func Test_Read_Queries_From_Files(t *testing.T) {
	// FIXME tests
	/*config := ReadConfig([]string{"Plugin1", "Plugin2"}, func(name string) io.ReadCloser {
		return ioutil.NopCloser(bytes.NewReader([]byte("all_exceptions=log:exception\nall_npe=log:NullPointerException")))
	})

	assert.Equal(t, len(config.byFile), 2)
	assert.NotEqual(t, config.ForPlugin("Plugin1"), nil)
	assert.NotEqual(t, config.ForPlugin("Plugin2"), nil)

	assert.Equal(t, config.ForPlugin("Plugin1")["all_exceptions"], "log:exception")
	assert.Equal(t, config.ForPlugin("Plugin1")["all_npe"], "log:NullPointerException")*/
}
