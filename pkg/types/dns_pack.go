package types

import "github.com/miekg/dns"

// DNSPack pack up with a dns message
type DNSPack struct {
	Msg *dns.Msg
}
