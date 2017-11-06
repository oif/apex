package hosts

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/oif/apex/pkg/config"
	plugin "github.com/oif/apex/pkg/plugin/v1"

	"github.com/Sirupsen/logrus"
	"github.com/miekg/dns"
)

// PluginName for plugin
const PluginName = "Hosts Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct {
	Config config.Hosts
	hosts  map[string]net.IP
	lock   sync.RWMutex
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Plugin
func (p *Plugin) Initialize() error {
	if p.Config.Enable {
		p.Config.Check() // check config

		go func() {
			for {
				startAt := time.Now().UnixNano()
				if err := p.hostsUpdater(); err != nil {
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Error("Update hosts failed")
				} else {
					logrus.WithFields(logrus.Fields{
						"time_usage":  time.Now().UnixNano() - startAt,
						"hosts_count": len(p.hosts),
					}).Debug("Update hosts success")
				}
				time.Sleep(time.Duration(p.Config.UpdateInterval) * time.Second)
			}
		}()
	}
	return nil
}

func (p *Plugin) hostsUpdater() error {
	newHosts, err := parseHosts(p.Config.Address)
	if err != nil {
		return err
	}

	// update hosts
	p.lock.Lock()
	p.hosts = newHosts
	p.lock.Unlock()

	return nil
}

func (p *Plugin) getHosts(name string) (net.IP, bool) {
	p.lock.RLock()
	ip, ok := p.hosts[name]
	p.lock.RUnlock()
	return ip, ok
}

func (p *Plugin) Warmup(c *plugin.Context)                   {}
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {}

func (p *Plugin) Patch(c *plugin.Context) {
	// for i, question := range c.Msg.Question {
	// 	if question.Qtype != dns.TypeA { // only hosts type A
	// 		continue
	// 	}
	// 	if c.Msg.Answer == nil {
	// 		log.Print("nil")
	// 	}
	// 	if ip, ok := p.getHosts(question.Name); ok { // match hosts
	// 		c.Msg.Answer[i] = &dns.A{
	// 			Hdr: dns.RR_Header{
	// 				Name:   question.Name,
	// 				Rrtype: dns.TypeA,
	// 				Ttl:    300,
	// 			},
	// 			A: ip,
	// 		}
	// 		c.Abort()
	// 		logrus.WithFields(logrus.Fields{
	// 			"name": question.Name,
	// 		}).Debug("Match hosts")
	// 	}
	// }

	if c.Msg.Question[0].Qtype != dns.TypeA { // only hosts type A
		return
	}
	if c.Msg.Answer == nil {
		log.Print("nil")
	}
	if ip, ok := p.getHosts(c.Msg.Question[0].Name); ok { // match hosts
		c.Msg.Answer = []dns.RR{
			&dns.A{
				Hdr: dns.RR_Header{
					Name:   c.Msg.Question[0].Name,
					Rrtype: dns.TypeA,
					Class:  c.Msg.Question[0].Qclass,
					Ttl:    300,
				},
				A: ip,
			},
		}
		c.Abort()
		logrus.WithFields(logrus.Fields{
			"name": c.Msg.Question[0].Name,
		}).Debug("Match hosts")
	}
}
