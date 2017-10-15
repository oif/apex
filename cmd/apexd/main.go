package main

import (
	"github.com/Sirupsen/logrus"
	inbound "github.com/oif/apex/pkg/inbound/v1"
	"github.com/oif/apex/plugin/gdns"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	s := new(inbound.Server)
	s.ListenAddress = ":53"
	s.ListenProtocol = []string{"udp"}
	s.RegisterPlugins(func() *gdns.Plugin { return new(gdns.Plugin) }())
	s.Run()
}
