package cache

import (
	plugin "github.com/oif/apex/pkg/plugin/v1"

	"github.com/coocood/freecache"
)

// PluginName for g.Name
const PluginName = "Cache Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct{}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Google DNS Plugin
func (p *Plugin) Initialize() error {
	cache = freecache.NewCache(512 * 1024) // kb
	return nil
}

// Warmup implements plugin
func (p *Plugin) Warmup(c *plugin.Context) {
	if err := getDNSCache(c.Msg, c.ClientIP().String()); err != nil { // miss cache
		c.Logger().Debug("Miss cache")
		c.Set("cache_plugin:write", true)
	} else { // get cache
		c.Logger().Debug("Hit cache")
		c.Abort()
	}
}

// AfterResponse implements plugin
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {
	if shouldWriteCache := c.GetBool("cache_plugin:write"); shouldWriteCache {
		c.Logger().Debug("Write cache")
		addDNSCache(c.Msg, c.ClientIP().String())
	}
}

// Patch the dns pakcage
func (p *Plugin) Patch(c *plugin.Context) {}
