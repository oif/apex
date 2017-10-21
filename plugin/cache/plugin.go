package cache

import (
	plugin "github.com/oif/apex/pkg/plugin/v1"

	lru "github.com/hashicorp/golang-lru"
)

// PluginName for g.Name
const PluginName = "Cache Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct {
	CacheSize int
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Google DNS Plugin
func (p *Plugin) Initialize() (err error) {
	cache, err = lru.New(p.CacheSize)
	return
}

// Warmup implements plugin
func (p *Plugin) Warmup(c *plugin.Context) {
	if hit := getCache(c.Msg, c.ClientIP()); hit { // hit cache
		c.Logger().Debug("Hit cache")
		c.Abort()
	} else { // miss cache
		c.Logger().Debug("Miss cache")
		c.Set("cache_plugin:write", true)
	}
}

// AfterResponse implements plugin
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {
	if shouldWriteCache := c.GetBool("cache_plugin:write"); shouldWriteCache {
		c.Logger().Debug("Write cache")
		writeCache(c.Msg, c.ClientIP())
	}
}

// Patch the dns pakcage
func (p *Plugin) Patch(c *plugin.Context) {}
