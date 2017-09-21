package types

// EDNSClientSubnet edns-client-subnet
type EDNSClientSubnet struct {
	Enable     string
	ExternalIP string // if external IP is empty will use request client IP
}
