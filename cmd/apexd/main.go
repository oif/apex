package main

import (
	"github.com/oif/apex/pkg/config"
	engine "github.com/oif/apex/pkg/engine/v1"
	"github.com/oif/apex/plugin/cache"
	"github.com/oif/apex/plugin/gdns"
	"github.com/oif/apex/plugin/hosts"
	"github.com/oif/apex/plugin/mupstream"
	"github.com/oif/apex/plugin/statistics"

	log "github.com/Sirupsen/logrus"
)

func main() {
	log.SetLevel(log.WarnLevel)
	s := new(engine.Server)
	s.ListenAddress = ":53"
	s.ListenProtocol = []string{"udp"}
	s.RegisterPlugins(func() *statistics.Plugin {
		plugin := new(statistics.Plugin)
		plugin.ConfigFilePath = "statistics.toml"
		return plugin
	}())
	s.RegisterPlugins(func() *hosts.Plugin {
		plugin := new(hosts.Plugin)
		plugin.Config = config.Hosts{
			Enable:         true,
			UpdateInterval: 3600,
			Address:        "https://raw.githubusercontent.com/racaljk/hosts/master/hosts",
		}
		return plugin
	}())
	s.RegisterPlugins(func() *cache.Plugin {
		plugin := new(cache.Plugin)
		plugin.CacheSize = 1024
		return plugin
	}())
	s.RegisterPlugins(func() *mupstream.Plugin {
		plugin := new(mupstream.Plugin)
		return plugin
	}())
	s.RegisterPlugins(func() *gdns.Plugin {
		plugin := new(gdns.Plugin)
		plugin.EnableProxy = true
		plugin.ProxyAddr = "http://127.0.0.1:8118"
		return plugin
	}())
	s.Run()
}
