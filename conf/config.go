package conf

import (
	"fmt"
	"github.com/caarlos0/env"
)

//ServerConfig represents Main service configuration
type ServerConfig struct {
	Hostname string `env:"HOSTNAME" envDefault:"localhost"`
	Port     int    `env:"RP_SERVER_PORT" envDefault:"8080"`
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
func EmptyConfig() *ServerConfig {
	return &ServerConfig{}
}
