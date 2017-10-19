package statistics

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	InfluxDB InfluxDB
}

type InfluxDB struct {
	Addr     string
	Username string
	Password string
}

// Load config from file
func (c *Config) Load(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file [%s] not found", filename))
	}
	if _, err := toml.DecodeFile(filename, c); err != nil {
		panic(fmt.Sprintf("Failed to decode config: %s", err))
	}
}
