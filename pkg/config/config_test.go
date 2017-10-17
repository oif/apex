package config

import (
	"testing"

	"github.com/oif/apex/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	config := NewConfig()
	assert.NotPanics(t, func() {
		config.Load("test.toml").Check()
	}, "This test case should not panic")
	assert.Equal(t, DefaultListenAddress, config.ListenAddress, "Set default listen address where ListenAddress is empty in config")
	assert.Equal(t, DefaultListenProtocol, config.ListenProtocol[0], "Set default listen protocol where ListenProtocol is empty in config")

	config = NewConfig()
	config.ListenProtocol = []string{"http"}
	assert.Panics(t, func() { config.Check() }, "Panic due to invalid listen protocol")

	config = NewConfig()
	assert.Panics(t, func() {
		config.Load("no_exists.toml")
	}, "Should panic due to file not exists")

	config = NewConfig()
	assert.Panics(t, func() {
		config.Load("test_bad.toml")
	}, "Should panic due to invalid toml format file given")
}

func TestUpstreamCheck(t *testing.T) {
	ups := &Upstream{}
	assert.Panics(t, func() {
		ups.Check()
	}, "Should panic due to primary DNS and alternative DNS are both empty")

	ups.PrimaryDNS = append(ups.PrimaryDNS, &types.Upstream{})
	assert.Panics(t, func() {
		ups.Check()
	}, "Should panic due to primary DNS check fail")

	ups.PrimaryDNS = make([]*types.Upstream, 0)

	ups.AlternativeDNS = append(ups.AlternativeDNS, &types.Upstream{})
	ups.AlternativeDNS = append(ups.PrimaryDNS, &types.Upstream{})
	assert.Panics(t, func() {
		ups.Check()
	}, "Should panic due to alternative DNS check fail")

	ups.AlternativeDNS[0] = &types.Upstream{
		Name:     "Google DNS",
		Address:  "8.8.8.8:53",
		Protocol: "udp",
	}
	assert.NotPanics(t, func() {
		ups.Check()
	}, "No panic here")
}

func TestProxyCheck(t *testing.T) {
	p := &Proxy{
		Policy: "",
	}

	assert.Panics(t, func() {
		p.Check()
	}, "Should Panic due to invalid policy")

	p.Policy = types.ProxyActivePolicy
	assert.Panics(t, func() {
		p.Check()
	}, "Should panic due to did not set proxy")

	p.Proxy = &types.Proxy{
		Protocol: "http",
	}
	assert.Panics(t, func() {
		p.Check()
	}, "Should panic due to invalid proxy")

	p.Proxy = &types.Proxy{
		Protocol: "http",
		Address:  "127.0.0.1:1080",
	}
	assert.NotPanics(t, func() {
		p.Check()
	}, "No panic here")
}

func TestCacheCheck(t *testing.T) {
	c := &Cache{
		Enable: true,
		Size:   -1,
	}

	assert.Panics(t, func() {
		c.Check()
	}, "Should panic due to invalid size")

	c.Size = 2
	assert.NotPanics(t, func() {
		c.Check()
	}, "No panic here")
}

func TestHostsCheck(t *testing.T) {
	h := &Hosts{
		Enable: true,
	}

	assert.Panics(t, func() {
		h.Check()
	}, "Should panic due to empty address")

	h.Address = "https://line.cat/hosts"
	assert.Panics(t, func() {
		h.Check()
	}, "Should panic due to update interval < 0")

	h.UpdateInterval = 3600
	assert.NotPanics(t, func() {
		h.Check()
	}, "No panic here")
}
