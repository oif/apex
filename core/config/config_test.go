package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := NewConfig()
	assert.NotPanics(t, func() {
		config.Load("test.toml").Check()
	}, "This test case should not panic")

	config = NewConfig()
	assert.Panics(t, func() {
		config.Load("no_exists.toml")
	}, "Should panic due to file not exists")
}
