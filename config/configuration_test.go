package config

import (
	"os"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	os.Setenv("GIGAPI_ROOT", "data")
	os.Setenv("GIGAPI_LAYERS_0_NAME", "fs")
	os.Setenv("GIGAPI_LAYERS_0_TYPE", "fs")
	os.Setenv("GIGAPI_LAYERS_0_GLOBAL", "false")
	os.Setenv("GIGAPI_LAYERS_0_URL", "file:///data/folder")

	os.Setenv("GIGAPI_LAYERS_1_NAME", "s3")
	os.Setenv("GIGAPI_LAYERS_1_TYPE", "s3")
	os.Setenv("GIGAPI_LAYERS_1_GLOBAL", "true")
	os.Setenv("GIGAPI_LAYERS_1_URL", "s3://key:secret@localhost:8181/bucket/prefix")

	InitConfig("")
	if Config.Gigapi.Root != "data" {
		t.Error("Expected GIGAPI_ROOT to be 'data', got", Config.Gigapi.Root)
	}
}
