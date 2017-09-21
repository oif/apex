package config

import "testing"

func TestConfig(t *testing.T) {
	config := new(Config)
	config.Load("test.toml")
}
