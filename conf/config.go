package conf

import (
	"fmt"
	"github.com/caarlos0/env"
)

//Registry represents type of used service discovery server
type Registry string

const (
	//Consul service discovery
	Consul Registry = "consul"

	// Eureka support has been dropped starting from version 1.1
	//Eureka service discovery
	//Eureka Registry = "eureka"
)

//ServerConfig represents Main service configuration
type ServerConfig struct {
	Hostname string `env:"HOSTNAME" envDefault:"localhost"`
	Port     int    `env:"RP_SERVER_PORT" envDefault:"8080"`
}

//ConsulConfig represents Consul Discovery service configuration
type ConsulConfig struct {
	Address      string   `env:"RP_CONSUL_ADDRESS" envDefault:"registry:8500"`
	Scheme       string   `env:"RP_CONSUL_SCHEME" envDefault:"http"`
	Token        string   `env:"RP_CONSUL_TOKEN"`
	PollInterval int      `env:"RP_CONSUL_POLL_INTERVAL" envDefault:"5"`
	PreferIP     bool     `env:"RP_CONSUL_PREFER_IP_ADDRESS" envDefault:"false"`
	Tags         []string `env:"RP_CONSUL_TAGS"`
}

//AddTags parses tags to string array. Extremely slow implementation - simplicity over speed
func (c *ConsulConfig) AddTags(tags ...string) {
	c.Tags = append(c.Tags, tags...)
}

//RpConfig represents Composite of all app configs
type RpConfig struct {
	AppName  string   `env:"RP_APP_NAME" envDefault:"goRP"`
	Registry Registry `env:"RP_REGISTRY" envDefault:"consul"`
	Server   *ServerConfig
	Consul   *ConsulConfig
}

//LoadConfig loads configuration from provided file and serializes it into RpConfig struct
func LoadConfig(cfg interface{}) error {
	err := env.Parse(cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}

	return nil
}

//EmptyConfig creates empty config
func EmptyConfig() *RpConfig {
	return &RpConfig{
		Consul: &ConsulConfig{},
		Server: &ServerConfig{},
	}
}
