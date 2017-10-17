package v1

// PluginChain contain series of plugins
type PluginChain []Object

// PluginFunc type
type PluginFunc func(*Context)

// PluginFuncChain type
type PluginFuncChain []PluginFunc
