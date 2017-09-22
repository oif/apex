package config

import "testing"

func TestConfig(t *testing.T) {
	config := NewConfig()
	config.Load("test.toml").Check()
}
