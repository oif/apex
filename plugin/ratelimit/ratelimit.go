package ratelimit

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	plugin "github.com/oif/apex/pkg/plugin/v1"
)

type clientLimiter struct {
	limiter    *rate.Limiter
	lastUpdate time.Time
}

var (
	RateLimit              = errors.New("hint rate limit")
	ReflectAttackSuspicion = errors.New("reflect attack suspicion")
)

func newClientLimiter() *clientLimiter {
	c := new(clientLimiter)
	c.limiter = rate.NewLimiter(rate.Every(time.Second), 20)
	c.lastUpdate = time.Now()
	return c
}

// PluginName for g.Name
const PluginName = "Rate Limit Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct {
	limiters      map[string]*clientLimiter
	reflectDefeat map[string]*ReflectDefeat
	lock          sync.RWMutex
}

type ReflectDefeat struct {
	IP           string
	RequestCount int
	RequestMap   map[string]bool
	LastRequest  time.Time
	Suspicion    bool
}

func newReflectDefeat(clientIP string) *ReflectDefeat {
	return &ReflectDefeat{
		IP:         clientIP,
		RequestMap: make(map[string]bool),
	}
}

func (r *ReflectDefeat) Record(domain string) {
	r.RequestMap[domain] = true
	r.RequestCount++
	r.LastRequest = time.Now()
}

func (r *ReflectDefeat) IsInactive() bool {
	// Inactive for 5min
	return time.Now().Sub(r.LastRequest).Seconds() > 3660
}

func (r *ReflectDefeat) HasSuspicion() bool {
	if r.Suspicion {
		return true
	}
	suspicion := r.RequestCount > 20 && len(r.RequestMap) < 2
	if suspicion {
		r.Suspicion = true
	}
	return r.Suspicion
}

func (r *ReflectDefeat) debug() string {
	return fmt.Sprintf("IP: %s, RequestCount: %d, RequestDomain: %d", r.IP, r.RequestCount, len(r.RequestMap))
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Google DNS Plugin
func (p *Plugin) Initialize() error {
	p.limiters = make(map[string]*clientLimiter)
	p.reflectDefeat = make(map[string]*ReflectDefeat)
	p.lock = sync.RWMutex{}

	go func() {
		for {
			now := time.Now()
			for k, l := range p.limiters {
				if now.Sub(l.lastUpdate) > 5*time.Minute {
					p.lock.Lock()
					delete(p.limiters, k)
					p.lock.Unlock()
				}
			}
			for k, r := range p.reflectDefeat {
				if r.IsInactive() {
					p.lock.Lock()
					delete(p.reflectDefeat, k)
					p.lock.Unlock()
				}
			}
			time.Sleep(5 * time.Minute)
		}
	}()
	return nil
}

// Warmup implements plugin
func (p *Plugin) Warmup(c *plugin.Context) {}

// AfterResponse implements plugin
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {}

// Patch the dns pakcage
func (p *Plugin) Patch(c *plugin.Context) {
	//if c.Msg.Question[0].Qtype != dns.TypeANY {
	//	return
	//}
	clientIP := c.ClientIP().String()
	reflectDefeat := p.getReflectDefeat(clientIP)
	//c.Logger().Infof("Reflect defeat -> %s", reflectDefeat.debug())

	if reflectDefeat.HasSuspicion() {
		c.Logger().WithField("clientIP", clientIP).Info("Suspicion request of reflect attack")
		c.AbortWithError(ReflectAttackSuspicion)
		return
	}
	reflectDefeat.Record(strings.ToLower(c.Msg.Question[0].Name))

	limiter := p.getLimiter(clientIP)
	if !limiter.limiter.Allow() {
		// resolve error
		c.Logger().WithField("clientIP", clientIP).Debug("Hit rate limit")
		c.AbortWithError(RateLimit)
		return
	}
}

func (p *Plugin) getReflectDefeat(clientIP string) *ReflectDefeat {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.reflectDefeat[clientIP]; !ok {
		p.reflectDefeat[clientIP] = newReflectDefeat(clientIP)
	} else {
		p.reflectDefeat[clientIP].LastRequest = time.Now()
	}
	return p.reflectDefeat[clientIP]
}

func (p *Plugin) getLimiter(clientIP string) *clientLimiter {
	p.lock.Lock()
	defer p.lock.Unlock()

	if _, ok := p.limiters[clientIP]; !ok {
		p.limiters[clientIP] = newClientLimiter()
	} else {
		p.limiters[clientIP].lastUpdate = time.Now()
	}
	return p.limiters[clientIP]
}
