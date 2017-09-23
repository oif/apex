package types

import "errors"

// Proxy type
type Proxy struct {
	Protocol string
	Address  string
}

const (
	// ProxyDisablePolicy never use proxy
	ProxyDisablePolicy = "disable"
	// ProxyActivePolicy use proxy always
	ProxyActivePolicy = "active"
	// ProxyPassivePolicy use in need
	ProxyPassivePolicy = "passive"
)

// Check proxy
func (p *Proxy) Check() error {
	if p.Protocol == "" {
		return errors.New("empty proxy protocol")
	}
	if p.Address == "" {
		return errors.New("empty proxy address")
	}
	return nil
}
