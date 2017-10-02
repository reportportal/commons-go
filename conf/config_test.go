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

	rpConf := EmptyConfig()
	err := LoadConfig(rpConf)
	Ω(err).ShouldNot(HaveOccurred())
	Expect(rpConf.Server.Hostname).ShouldNot(BeEmpty())
}

func TestLoadConfigWithParameters(t *testing.T) {
	os.Setenv("RP_PARAMETERS_PARAM", "env_value")

	rpConf := struct {
		*RpConfig
		Param string `env:"RP_PARAMETERS_PARAM"`
	}{RpConfig: EmptyConfig()}

	err := LoadConfig(&rpConf)
	Ω(err).ShouldNot(HaveOccurred())

	if "env_value" != rpConf.Param {
		t.Error("Config parser fails")
	}
}

func TestLoadConfigNonExisting(t *testing.T) {
	rpConf := EmptyConfig()
	err := LoadConfig(rpConf)
	Ω(err).ShouldNot(HaveOccurred())

	if 8080 != rpConf.Server.Port {
		t.Error("Should not return empty string for default config")
	}
}

func TestLoadConfigIncorrectFormat(t *testing.T) {
	rpConf := EmptyConfig()
	err := LoadConfig(rpConf)
	Ω(err).ShouldNot(HaveOccurred())

	if 8080 != rpConf.Server.Port {
		t.Error("Should return empty string for default config")
	}
}

func TestLoadStringArray(t *testing.T) {
	os.Setenv("RP_CONSUL_TAGS", "tag1,tag2,tag3")
	rpConf := EmptyConfig()
	err := LoadConfig(rpConf)
	Ω(err).ShouldNot(HaveOccurred())

	rpConf.Consul.AddTags("tag4", "tag5")
	fmt.Println(rpConf.Consul.Tags)
	expected := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	if !reflect.DeepEqual(expected, rpConf.Consul.Tags) {
		t.Errorf("Incorrect array parameters parsing. Expected: %s, Actual: %s", expected, rpConf.Consul.Tags)
	}
}
