package v1

import (
	"sync"
	"time"

	plugin "github.com/oif/apex/pkg/plugin/v1"

	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
	"github.com/sony/sonyflake"
)

// Server implements DNS server with dns.Handler
type Server struct {
	ListenAddress  string
	ListenProtocol []string

	mux  *dns.ServeMux
	srvs []*dns.Server
	lock sync.RWMutex
	wg   *sync.WaitGroup
	uuid *sonyflake.Sonyflake

	// plugin
	plugins plugin.PluginChain
}

// Run server
func (s *Server) Run() {
	log.Infof("Starting service at %v on %v", s.ListenAddress, s.ListenProtocol)

	s.mux = dns.NewServeMux()
	log.Debug("Setting DNS server handle")
	s.mux.Handle(".", s)

	s.wg = new(sync.WaitGroup)
	s.uuid = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return 1, nil
		},
		CheckMachineID: func(id uint16) bool {
			return true
		},
	})
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
	var (
		// abort bool
		err     error
		reqID   uint64
		context *plugin.Context
	)
	reqID, err = s.uuid.NextID()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Get request ID failed")
		goto RESPONSE
	}

	log.WithFields(log.Fields{
		"req_id": reqID,
		"name":   m.Question[0].Name,
		"type":   m.Question[0].Qtype,
	}).Debug("Receive request")

	context = plugin.NewContext(w, m, reqID)
	context.MustRegisterPluginsOnce(s.plugins)

	context.Warmup()
	context.Patch()

	if context.HasError() {
		for _, err = range context.Errors {
			log.Errorf("Context error: %v", err)
		}
	}

	log.WithFields(log.Fields{
		"req_id": reqID,
	}).Debug("Resolve done ready to response")

RESPONSE:
	// write resposne message
	if err = w.WriteMsg(m); err != nil {
		log.Errorf("Error when write response message: %v", err)
	}
	context.AfterResponse(err)
}

// RegisterPlugins for server
func (s *Server) RegisterPlugins(p plugin.Object) error {
	s.plugins = append(s.plugins, p)
	// @TODO do some initialization works here
	return p.Initialize()
}

// Stop the server graceful
func (s *Server) Stop() {
	s.lock.Lock()
	defer s.lock.Unlock()

	// @TODO should have a graceful way to close dns server and clean up
	log.Infoln("Graceful shutdown")
	for _, srv := range s.srvs {
		srv.Shutdown()
	}
	s.wg.Done()
}
