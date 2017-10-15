package v1

import (
	"math"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

const (
	abortIndex int8 = math.MaxInt8 / 2 // 63
)

// Context pack up with a dns message, and is the most important part of apex
// It allow us to pass variables throught plugins, manage the work flow
type Context struct {
	Writer dns.ResponseWriter
	Msg    *dns.Msg
	Errors []error
	KV     map[string]interface{}

	clientIP  net.IP
	requestID uint64

	// plugin
	index         int8
	plugins       PluginChain // @TODO postStart preStop ?
	pluginsLength int8
}

// NewContext with basic properties
func NewContext(w dns.ResponseWriter, m *dns.Msg, reqID uint64) *Context {
	c := new(Context)
	c.Writer = w
	c.Msg = m
	c.requestID = reqID
	return c
}

// MustRegisterPluginsOnce register plugins
func (c *Context) MustRegisterPluginsOnce(pluginsChain PluginChain) {
	if pluginsChain == nil {
		panic("Plugin chain is nil")
	}
	c.plugins = pluginsChain
	c.pluginsLength = int8(len(c.plugins))
}

// Next plugin
func (c *Context) Next() {
	for ; c.index < c.pluginsLength; c.index++ {
		c.plugins[c.index].Patch(c)
	}
}

// Abort prevents plugins calling after current plugin
func (c *Context) Abort() {
	c.index = abortIndex
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Error will panic if err is nil. Append err into context.errors
func (c *Context) Error(err error) error {
	if err == nil {
		panic("except error but nil given")
	}
	c.Errors = append(c.Errors, err)
	return err
}

// AbortWithError calls `Abort()` and `Error()`
func (c *Context) AbortWithError(err error) error {
	c.Abort()
	return c.Error(err)
}

// HasError returns true if the current context has some errors.
func (c *Context) HasError() bool {
	return len(c.Errors) > 0
}

// Logger returns logrus.Entry with request ID
func (c *Context) Logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"req_id": c.requestID,
	})
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
func (c *Context) ClientIP() net.IP {
	if c.clientIP == nil {
		switch c.Writer.RemoteAddr().Network() {
		case "tcp", "tcp4", "tcp6":
			c.clientIP = c.Writer.RemoteAddr().(*net.TCPAddr).IP
		case "udp", "udp4", "udp6":
			c.clientIP = c.Writer.RemoteAddr().(*net.UDPAddr).IP
		}

		if !isInternalIP(c.clientIP) {
			// get dns internal ip
		}
	}
	return c.clientIP
}

// Deadline implements context
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done implements context
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err implements context
func (c *Context) Err() error {
	return nil
}

// Value implements context
func (c *Context) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
