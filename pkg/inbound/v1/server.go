package v1

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

// Server implements DNS server with dns.Handler
type Server struct {
	ListenAddress  string
	ListenProtocol []string

	mux  *dns.ServeMux
	srvs []*dns.Server
	lock sync.RWMutex
	wg   *sync.WaitGroup
}

// Run server
func (s *Server) Run() {
	log.Infof("Starting service at %s on ", s.ListenAddress, s.ListenProtocol)

	s.mux = dns.NewServeMux()
	log.Debugln("Setting DNS server handle")
	s.mux.Handle(".", s)

	s.wg = new(sync.WaitGroup)
	// Add wait group
	s.wg.Add(len(s.ListenProtocol))

	for _, proto := range s.ListenProtocol {
		log.Infof("Ready to serve at [%s]%s", proto, s.ListenAddress)
		go func(network string) {
			s.lock.Lock()
			log.Debug("Server locked")
			srv := &dns.Server{
				Addr:    s.ListenAddress,
				Net:     network,
				Handler: s.mux,
			}
			s.srvs = append(s.srvs, srv)
			s.lock.Unlock()
			log.Debug("Server unlocked")

			if err := srv.ListenAndServe(); err != nil {
				log.Panicf("Serve at [%s]%s error: %v", network, s.ListenAddress, err)
			}
		}(proto)
	}
	s.wg.Wait()
}

// ServeDNS implements dns.Handler interface
func (s *Server) ServeDNS(w dns.ResponseWriter, m *dns.Msg) {
	log.Debugf("Receive request\n%v", m)
}

// Stop the server graceful
func (s *Server) Stop() {
	// @TODO should have a graceful way to close dns server and clean up
	log.Infoln("Graceful shutdown")
	for _, srv := range s.srvs {
		srv.Shutdown()
	}
	s.wg.Done()
}
