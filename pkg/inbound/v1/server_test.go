package v1

import (
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	s := &Server{
		ListenAddress:  ":53",
		ListenProtocol: []string{"udp", "tcp"},
	}
	if os.Getenv("TRAVIS") != "true" {
		s.Run()
	}
}
