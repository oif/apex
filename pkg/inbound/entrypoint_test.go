package inbound

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

func TestEntrypoint(t *testing.T) {
	entrypoint := new(Entrypoint)
	entrypoint.HandleFunc = func(w dns.ResponseWriter, m *dns.Msg) {
		fmt.Printf("resolve request body %v\n", m)
		// no reply yet
	}
	entrypoint.Setup(&dns.Server{
		Addr: ":53",
		Net:  "udp",
	}).Serve()
}