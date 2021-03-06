package main

import (
	engine "github.com/oif/apex/pkg/engine/v1"
	"github.com/oif/apex/plugin/cache"
	"github.com/oif/apex/plugin/gdns"
	"github.com/oif/apex/plugin/mupstream"
	"github.com/oif/apex/plugin/ratelimit"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.InfoLevel)
	s := new(engine.Server)
	s.ListenAddress = ":53"
	s.ListenProtocol = []string{"udp"}
	s.RegisterPlugins(func() *ratelimit.Plugin {
		plugin := new(ratelimit.Plugin)
		return plugin
	}())
	s.RegisterPlugins(func() *cache.Plugin {
		plugin := new(cache.Plugin)
		plugin.CacheSize = 1024
		return plugin
	}())
	s.RegisterPlugins(func() *gdns.Plugin {
		plugin := new(gdns.Plugin)
		plugin.EnableProxy = true
		plugin.ProxyAddr = "http://127.0.0.1:8118"
		return plugin
	}())
	s.RegisterPlugins(func() *mupstream.Plugin {
		plugin := new(mupstream.Plugin)
		return plugin
	}())
	s.Run()
}
