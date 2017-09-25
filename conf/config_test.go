package conf

import (
	"fmt"
	. "github.com/onsi/gomega"
	"os"
	"reflect"
	"testing"
)

func TestLoadEmptyConfig(t *testing.T) {
	RegisterTestingT(t)

	rpConf, err := LoadConfig(nil, nil)
	Ω(err).ShouldNot(HaveOccurred())
	Expect(rpConf.Server.Hostname).ShouldNot(BeEmpty())
}

func TestLoadConfigWithParameters(t *testing.T) {
	os.Setenv("RP_PARAMETERS_PARAM", "env_value")
	rpConf, err := LoadConfig(EmptyConfig(), nil)
	Ω(err).ShouldNot(HaveOccurred())

	if "env_value" != rpConf.Get("RP_PARAMETERS_PARAM") {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigNonExisting(t *testing.T) {
	rpConf, err := LoadConfig(EmptyConfig(), nil)
	Ω(err).ShouldNot(HaveOccurred())

	if 8080 != rpConf.Server.Port {
		t.Error("Should not return empty string for default config")
	}
}

func TestLoadConfigIncorrectFormat(t *testing.T) {
	rpConf, err := LoadConfig(EmptyConfig(), nil)
	Ω(err).ShouldNot(HaveOccurred())

	if 8080 != rpConf.Server.Port {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadStringArray(t *testing.T) {
	os.Setenv("RP_CONSUL_TAGS", "tag1,tag2,tag3")
	rpConf, err := LoadConfig(nil, nil)
	Ω(err).ShouldNot(HaveOccurred())

	rpConf.Consul.AddTags("tag4", "tag5")
	fmt.Println(rpConf.Consul.Tags)
	expected := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	if !reflect.DeepEqual(expected, rpConf.Consul.Tags) {
		t.Errorf("Incorrect array parameters parsing. Expected: %s, Actual: %s", expected, rpConf.Consul.Tags)
	}
}
