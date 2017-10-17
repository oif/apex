package config

import (
	"fmt"
	"os"

	"github.com/oif/apex/pkg/types"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultListenAddress set :53
	DefaultListenAddress = ":53"
	// DefaultListenProtocol udp
	DefaultListenProtocol = ProtocolUDP
	// ProtocolUDP constant
	ProtocolUDP = "udp"
	// ProtocolTCP constant
	ProtocolTCP = "tcp"
)

// Config groups
type Config struct {
	ListenAddress  string   // default listen at :dns
	ListenProtocol []string // udp tcp only, default udp
	Upstream       *Upstream
	Proxy          *Proxy
	Cache          *Cache
	Hosts          *Hosts
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

// NewConfig return
func NewConfig() *Config {
	return &Config{
		Upstream: new(Upstream),
		Proxy: &Proxy{
			Policy: types.ProxyDisablePolicy,
		},
		Cache: new(Cache),
		Hosts: new(Hosts),
	}
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
	if c.ListenAddress == "" {
		c.ListenAddress = DefaultListenAddress
	}
	if len(c.ListenProtocol) == 0 {
		c.ListenProtocol = append(c.ListenProtocol, DefaultListenProtocol)
	} else {
		for _, protocol := range c.ListenProtocol {
			if protocol != ProtocolTCP && protocol != ProtocolUDP {
				panic(fmt.Sprintf("Unsupport DNS server listen protocol: %s", protocol))
			}
		}
	}

	if c.Upstream != nil {
		c.Upstream.Check()
	}
	if c.Proxy != nil {
		c.Proxy.Check()
	}
	if c.Cache != nil {
		c.Cache.Check()
	}
	if c.Hosts != nil {
		c.Hosts.Check()
	}
}

// Check upstream config
func (u *Upstream) Check() {
	if len(u.PrimaryDNS) == 0 && len(u.AlternativeDNS) == 0 {
		panic(fmt.Sprintf("Primary DNS and Alternative DNS could not be empty both, Primary DNS %d, Alternative DNS %d", len(u.PrimaryDNS), len(u.AlternativeDNS)))
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
	if p.Policy != types.ProxyActivePolicy && p.Policy != types.ProxyPassivePolicy && p.Policy != types.ProxyDisablePolicy {
		panic(fmt.Sprintf("Invalid proxy policy: %s", p.Policy))
	}
	if p.Policy != types.ProxyDisablePolicy {
		if p.Proxy == nil {
			panic(fmt.Sprintf("Use %s proxy policy but did not set proxy", p.Policy))
		}
		if err := p.Proxy.Check(); err != nil {
			panic(fmt.Sprintf("Invalid proxy : %s", err))
		}
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
