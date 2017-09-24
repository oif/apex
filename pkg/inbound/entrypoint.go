package inbound

import (
	"errors"

	"github.com/miekg/dns"
)

// Entrypoint of inbound
// @TODO plugin support
type Entrypoint struct {
	HandleFunc func(dns.ResponseWriter, *dns.Msg)
	Server     *dns.Server
	err        error
}

func (e *Entrypoint) hasError() error {
	if e == nil {
		return errors.New("this entrypoint is nil")
	}
	if e.Server == nil {
		return errors.New("dns server did not initialize, plz user Setup func")
	}
	if e.err != nil {
		return e.err
	}
	return nil
}

// Setup DNS server with config, where is the start func of entrypoint
func (e *Entrypoint) Setup(server *dns.Server) *Entrypoint {
	if server == nil {
		e.err = errors.New("DNS server is required")
		return e
	}
	e.Server = server
	if e.HandleFunc == nil {
		e.err = errors.New("handle func cloud not be nil")
		return e
	}
	dns.HandleFunc(".", e.HandleFunc)
	return e
}

// Serve run a DNS server
func (e *Entrypoint) Serve() {
	if err := e.hasError(); err != nil {
		panic(err)
	}

	if err := e.Server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Shutdown entrypoint
func (e *Entrypoint) Shutdown() error {
	if err := e.hasError(); err != nil {
		return err
	}
	// Do other clean works, such as close all the conn to do a graceful-shutdown
	return e.Server.Shutdown()
}
