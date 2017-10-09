package v1

import (
	"math"

	"github.com/miekg/dns"
)

const (
	abortIndex int8 = math.MaxInt8 / 2 // 63
)

// Context pack up with a dns message, and is the most important part of apex
// It allow us to pass variables throught plugins, manage the lifecycle and flow
type Context struct {
	Writer dns.ResponseWriter
	Msg    *dns.Msg
	Errors []error
	KV     map[string]interface{}

	index         int8
	plugins       PluginChain
	pluginsLength int8
}

// MustRegisterPluginsOnce register plugins
func (c *Context) MustRegisterPluginsOnce(objs ...Object) {
	c.plugins = objs
	c.pluginsLength = int8(len(c.plugins))
}

// Next plugin
func (c *Context) Next() {
	c.index++
	for ; c.index < c.pluginsLength; c.index++ {
		c.plugins[c.index].Patch(c)
	}
}

// Key Value Pair

// Set key value pair into context
func (c *Context) Set(key string, value interface{}) {
	if c.KV == nil {
		c.KV = make(map[string]interface{})
	}
	c.KV[key] = value
}

// Get value by key
func (c *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = c.KV[key]
	return
}

// GetString by key
func (c *Context) GetString(key string) (value string) {
	if temp, ok := c.KV[key]; ok {
		value = temp.(string)
	}
	return
}

// GetInt by key
func (c *Context) GetInt(key string) (value int) {
	if temp, ok := c.KV[key]; ok {
		value = temp.(int)
	}
	return
}

// GetBool by key
func (c *Context) GetBool(key string) (value bool) {
	if temp, ok := c.KV[key]; ok {
		value = temp.(bool)
	}
	return
}

// GetInt64 by key
func (c *Context) GetInt64(key string) (value int64) {
	if temp, ok := c.KV[key]; ok {
		value = temp.(int64)
	}
	return
}

// GetFloat64 by key
func (c *Context) GetFloat64(key string) (value float64) {
	if temp, ok := c.KV[key]; ok {
		value = temp.(float64)
	}
	return
}

// ClientIP return client request IP
func (c *Context) ClientIP() {}
