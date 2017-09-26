package v1

import "testing"

func TestServer(t *testing.T) {
	s := &Server{
		ListenAddress:  ":53",
		ListenProtocol: []string{"udp", "tcp"},
	}
	s.Run()
}
