package v1

// temporarily commented due to some issue of dns graceful shutdown https://github.com/miekg/dns/issues/457
// func TestServer(t *testing.T) {
// 	s := &Server{
// 		ListenAddress:  ":53",
// 		ListenProtocol: []string{"udp", "tcp"},
// 	}
// 	go func() {
// 		time.Sleep(10 * time.Second)
// 		s.Stop()
// 	}()
// 	s.Run()
// }
