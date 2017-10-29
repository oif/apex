package mupstream

import (
	"github.com/Sirupsen/logrus"
	plugin "github.com/oif/apex/pkg/plugin/v1"
)

// PluginName for plugin
const PluginName = "Multi-upstream Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct {
	upstreams []*upstream
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Google DNS Plugin
func (p *Plugin) Initialize() error {
	p.upstreams = append(p.upstreams, newUpstream("119.29.29.29:53"))
	p.upstreams = append(p.upstreams, newUpstream("223.5.5.5:53"))
	p.upstreams = append(p.upstreams, newUpstream("114.114.114.114:53"))
	return nil
}

func (p *Plugin) Warmup(c *plugin.Context)                   {}
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {}

func (p *Plugin) Patch(c *plugin.Context) {
	up := p.bestUpstream()
	resp, rtt, err := up.forward(c.Msg)
	if err != nil {
		c.Error(err)
		return
	}
	resp.CopyTo(c.Msg)
	c.Logger().WithFields(logrus.Fields{
		"rtt":      rtt,
		"upstream": up,
	}).Debug("Exchange message")

	c.Abort() // stop other patch steps
}

func (p *Plugin) bestUpstream() *upstream {
	best := 0
	for i := 0; i < len(p.upstreams); i++ {
		if p.upstreams[i].srtt < p.upstreams[0].srtt {
			best = i
		}
	}
	go func(selected int) { // lost decay
		for i := 0; i < len(p.upstreams); i++ {
			if i != selected {
				p.upstreams[i].srttAttenuation()
			}
		}
	}(best)

	return p.upstreams[best]
}
