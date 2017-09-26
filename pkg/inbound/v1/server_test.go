package v1

import (
	"os"
	"testing"
)

func TestServer(t *testing.T) {
	if os.Getenv("TRAVIS") != "TRUE" {
		s := &Server{
			ListenAddress:  ":53",
			ListenProtocol: []string{"udp", "tcp"},
		}
		s.Run()
	}
}
