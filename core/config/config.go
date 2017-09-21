package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/oif/apex/core/types"
)

// Config groups
type Config struct {
	ListenAddress string
	Upstream      *Upstream
	Proxy         *Proxy
	Cache         *Cache
	Hosts         *Hosts
}

// Upstream config
type Upstream struct {
	PrimaryDNS     []*types.Upstream
	AlternativeDNS []*types.Upstream
}

// Proxy config
type Proxy struct {
	Policy string
	Proxy  *types.Proxy
}

// Cache config
type Cache struct {
	Enable bool
	Size   int
}

// Hosts config
type Hosts struct {
	Enable         bool
	UpdateInterval int // update interval, unit second
	Address        string
}

// Load config from file
func (c *Config) Load(filename string) *Config {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		panic(fmt.Sprintf("Config file [%s] not found", filename))
	}
	if _, err := toml.DecodeFile(filename, c); err != nil {
		panic(fmt.Sprintf("Failed to decode config: %s", err))
	}
	return c
}

// Check values in config
func (c *Config) Check() {
	c.Upstream.Check()
	c.Proxy.Check()
	c.Cache.Check()
	c.Hosts.Check()
}

// Check upstream config
func (u *Upstream) Check() {
	if len(u.PrimaryDNS) == 0 && len(u.AlternativeDNS) == 0 {
		panic(fmt.Sprintf("Primary DNS and Alternative DNS could not be empty both"))
	}
	for i := 0; i < len(u.PrimaryDNS); i++ {
		if err := u.PrimaryDNS[i].Check(); err != nil {
			panic(fmt.Sprintf("Primary DNS[%d] %s", i, err))
		}
	}
	for i := 0; i < len(u.AlternativeDNS); i++ {
		if err := u.AlternativeDNS[i].Check(); err != nil {
			panic(fmt.Sprintf("Alternative DNS[%d] %s", i, err))
		}
	}
}

// Check proxy config
func (p *Proxy) Check() {
	if p.Policy != types.ProxyActivePolicy && p.Policy != types.ProxyPassivePolicy {
		panic(fmt.Sprintf("Invalid proxy policy: %s", p.Policy))
	}
	if err := p.Proxy.Check(); err != nil {
		panic(fmt.Sprintf("Invalid proxy : %s", err))
	}
}

// Check cache config
func (c *Cache) Check() {
	if c.Enable && c.Size < 0 {
		panic(fmt.Sprintf("Invalid cache size: %d", c.Size))
	}
}

// Check hosts config
func (h *Hosts) Check() {
	if h.Enable {
		if h.Address == "" {
			panic(fmt.Sprintf("Invalid address: %s", h.Address))
		}
		if h.UpdateInterval < 1 {
			panic(fmt.Sprintf("Hosts update interval should grater than 0, but %d given", h.UpdateInterval))
		}
	}
}
