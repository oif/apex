package v1

// Object of Plugin system for APEX inject into lifecycle of DNS server handle
type Object interface {
	// Name get plugin object name
	Name() string
	// Initialize plugin object
	Initialize() error

	// Warmup before patch
	Warmup(*Context)
	// Patch dns message and return it back for chain-call
	Patch(*Context)
	// AfterResponse call after response with response status
	AfterResponse(*Context, error)
}

// PluginChain contain series of plugins
type PluginChain []Object

// PluginFunc type
type PluginFunc func(*Context)

// PluginFuncChain type
type PluginFuncChain []PluginFunc
