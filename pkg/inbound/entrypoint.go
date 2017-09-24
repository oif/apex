package inbound

import (
	"github.com/miekg/dns"
	"github.com/oif/apex/pkg/config"
)

// Entrypoint of inbound
// @TODO plugin support
type Entrypoint struct {
	Server *dns.Server
}

// Setup DNS server with config
func (e *Entrypoint) Setup(c *config.Config, handleFunc func(dns.ResponseWriter, *dns.Msg)) (err error) {
	e.Server = &dns.Server{
		Addr: c.ListenAddress,
		Net:  "udp",
	}
	dns.HandleFunc(".", handleFunc)
	return
}

// Serve run a DNS server
func (e *Entrypoint) Serve() {
	if e.Server == nil {
		panic("Entrypoint did not setup server")
	}
	if err := e.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Shutdown entrypoint
func (e *Entrypoint) Shutdown() error {
	// Do other clean works, such as close all the conn to do a graceful-shutdown
	return e.Server.Shutdown()
}
