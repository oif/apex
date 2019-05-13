package v1

import (
	"strings"
	"sync"
	"time"

	plugin "github.com/oif/apex/pkg/plugin/v1"
	"github.com/oif/apex/plugin/ratelimit"

	log "github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
	"github.com/sony/sonyflake"
)

var (
	DomainBlacklist = []string{
		"access-board.gov",
		"usadf.gov",
		"aids.gov",
		"efps.gov",
		"nel.gov",
		"nccih.nih.gov",
		"commerce.gov",
	}
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
		err         error
		reqID       uint64
		context     *plugin.Context
		embedLogger = log.WithFields(log.Fields{
			"req_id": reqID,
			"name":   m.Question[0].Name,
			"type":   m.Question[0].Qtype,
		})
	)
	startAt := time.Now()
	reqID, err = s.uuid.NextID()
	if err != nil {
		embedLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Get request ID failed")
		goto RESPONSE
	}

	context = plugin.NewContext(w, m, reqID)
	context.MustRegisterPluginsOnce(s.plugins)

	embedLogger = embedLogger.WithField("client_ip", context.ClientIP())

	embedLogger.Debug("Receive request")

	// Blacklist
	for _, domain := range DomainBlacklist {
		if strings.Contains(strings.ToLower(m.Question[0].Name), domain) {
			embedLogger.Debug("Trigger domain blacklist, drop")
			return
		}
	}

	context.Warmup()
	context.Patch()

	if context.HasError() {
		for _, err = range context.Errors {
			if err == ratelimit.ReflectAttackSuspicion || err == ratelimit.RateLimit {
				embedLogger.WithField("dropReason", err).Info("Package dropped")
				return
			}
			embedLogger.Errorf("Context error: %v", err)
		}
	}

	embedLogger.Debug("Resolve done ready to response")

RESPONSE:
	// write response message
	if err = w.WriteMsg(m); err != nil {
		embedLogger.Errorf("Error when write response message: %v", err)
	}
	context.AfterResponse(err)
	fields := log.Fields{
		"timeUsed": time.Now().Sub(startAt).String(),
	}
	if err := context.Err(); err == nil {
		fields["name"] = m.Question[0].Name
		fields["type"] = m.Question[0].Qtype
	} else {
		fields["error"] = err
	}
	embedLogger.WithFields(fields).Info("Resolved request")
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
