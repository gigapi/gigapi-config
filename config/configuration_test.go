package config

import (
	"os"
	"testing"
)

func TestEnvConfig(t *testing.T) {
	os.Setenv("GIGAPI_ROOT", "data")
	InitConfig("")
	if Config.Gigapi.Root != "data" {
		t.Error("Expected GIGAPI_ROOT to be 'data', got", Config.Gigapi.Root)
	}
}
