package adapter

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	assert "gopkg.in/go-playground/assert.v1"
)

func Test_Read_Queries_From_Files(t *testing.T) {

	provider := &QueryProvider{}
	provider.openQueryFile = func(name string) io.ReadCloser {
		return ioutil.NopCloser(bytes.NewReader([]byte("all_exceptions=log:exception\nall_npe=log:NullPointerException")))
	}

	provider.read("defaultFile", []string{"Plugin1", "Plugin2"})

	assert.Equal(t, len(provider.QuerySets), 3)
	assert.NotEqual(t, provider.QuerySets["default"], nil)
	assert.NotEqual(t, provider.QuerySets["Plugin1"], nil)
	assert.NotEqual(t, provider.QuerySets["Plugin2"], nil)

	assert.Equal(t, provider.QuerySets["Plugin1"].Queries["all_exceptions"], "log:exception")
	assert.Equal(t, provider.QuerySets["Plugin1"].Queries["all_npe"], "log:NullPointerException")
}
