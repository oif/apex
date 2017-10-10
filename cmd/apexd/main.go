package main

import (
	inbound "github.com/oif/apex/pkg/inbound/v1"
	"github.com/oif/apex/plugin"
)

func main() {
	s := new(inbound.Server)
	s.ListenAddress = ":53"
	s.ListenProtocol = []string{"udp"}
	s.RegisterPlugins(func() *plugin.GoogleDNS { return new(plugin.GoogleDNS) }())
	s.Run()
}
