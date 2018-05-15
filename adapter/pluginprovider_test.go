package adapter

import (
	"testing"

	assert "gopkg.in/go-playground/assert.v1"
)

func Test_Get_Logical_Plugin_Names(t *testing.T) {
	file := "path/file.so"
	expectedFile := "file"

	result := getLogicalPluginName(file)

	assert.Equal(t, result, expectedFile)
}
