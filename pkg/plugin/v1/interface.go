package v1

// Object of Plugin system for APEX inject into lifecycle of DNS server handle
type Object interface {
	// Name get plugin object name
	Name() string
	// Initialize plugin object
	Initialize() error
	// Patch dns message and return it back for chain-call
	Patch(*Context)
}
