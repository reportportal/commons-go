package conf

import (
	"os"
	"testing"
	"fmt"
	"reflect"
)

func TestLoadConfig(t *testing.T) {
	rpConf := LoadConfig("./../server.yaml", nil)
	if "10.200.10.1" != rpConf.Server.Hostname {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigWithParameters(t *testing.T) {
	os.Setenv("RP_PARAMETERS.PARAM", "env_value")
	rpConf := LoadConfig("", map[string]interface{}{"parameters.param": "default_value"})

	if "env_value" != rpConf.Get("parameters.param").(string) {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigNonExisting(t *testing.T) {
	rpConf := LoadConfig("server.yaml", nil)
	if 8080 != rpConf.Server.Port {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadConfigIncorrectFormat(t *testing.T) {
	rpConf := LoadConfig("config_test.go", nil)
	if 8080 != rpConf.Server.Port {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadStringArray(t *testing.T) {
	os.Setenv("RP_CONSUL.TAGS", "tag1,tag2,tag3")
	rpConf := LoadConfig("", nil)
	rpConf.Consul.AddTags("tag4", "tag5")
	fmt.Println(rpConf.Consul.Tags)
	expected := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	if !reflect.DeepEqual(expected, rpConf.Consul.GetTags()) {
		t.Errorf("Incorrect array parameters parsing. Expected: %s, Actual: %s", expected, rpConf.Consul.Tags)
	}
}
