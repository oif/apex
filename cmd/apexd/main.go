package main

import (
	inbound "github.com/oif/apex/pkg/engine/v1"
	"github.com/oif/apex/plugin/cache"
	"github.com/oif/apex/plugin/gdns"
	"github.com/oif/apex/plugin/statistics"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	s := new(inbound.Server)
	s.ListenAddress = ":53"
	s.ListenProtocol = []string{"udp"}
	s.RegisterPlugins(func() *statistics.Plugin {
		plugin := new(statistics.Plugin)
		plugin.ConfigFilePath = "statistics.toml"
		return plugin
	}())
	s.RegisterPlugins(func() *cache.Plugin { return new(cache.Plugin) }())
	s.RegisterPlugins(func() *gdns.Plugin { return new(gdns.Plugin) }())
	s.Run()
}
