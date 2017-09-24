package inbound

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
	"github.com/oif/apex/pkg/config"
)

func TestEntrypoint(t *testing.T) {
	entrypoint := new(Entrypoint)
	entrypoint.Setup(&config.Config{
		ListenAddress: ":53",
	}, func(w dns.ResponseWriter, m *dns.Msg) {
		fmt.Printf("resolve request body %v\n", m)
		// no reply yet
	})
	entrypoint.Serve()
}
