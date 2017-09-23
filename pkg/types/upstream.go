package types

import "errors"

// Upstream of DNS
type Upstream struct {
	Name             string
	Address          string // include IP address and port
	Protocol         string // udp, tcp, http, https. if is http or https means HttpDNS but current support Google DNS only
	Timeout          int8   // unit second, default timeout 3
	EnableProxy      bool   // default false
	EDNSClientSubnet EDNSClientSubnet
}

// DefaultUpstreamTimeout 3 seconds
const DefaultUpstreamTimeout = 3

// Check upstream
func (u *Upstream) Check() error {
	if u.Address == "" {
		return errors.New("DNS address is empty")
	}
	if u.Protocol == "" {
		return errors.New("DNS protocol is empty")
	}
	if u.Name == "" {
		return errors.New("DNS name is empty")
	}
	if u.Timeout < 0 {
		return errors.New("DNS timeout is invalid")
	} else if u.Timeout == 0 {
		u.Timeout = DefaultUpstreamTimeout
	}
	return nil
}
