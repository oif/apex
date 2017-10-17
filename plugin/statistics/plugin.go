package statistics

import (
	"time"

	"github.com/Sirupsen/logrus"
	plugin "github.com/oif/apex/pkg/plugin/v1"
)

// PluginName for g.Name
const PluginName = "Statistics Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct{}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Google DNS Plugin
func (p *Plugin) Initialize() error {
	return nil
}

// Warmup implements plugin
func (p *Plugin) Warmup(c *plugin.Context) {
	c.Set("statistics_plugin:startTime", makeTimestamp())
}

// AfterResponse implements plugin
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {
	if startAt := c.GetInt64("statistics_plugin:startTime"); startAt != 0 {
		c.Logger().WithFields(logrus.Fields{
			"response_time": makeTimestamp() - startAt,
		}).Info("Response time usage statistics")
	}
}

// Patch the dns pakcage
func (p *Plugin) Patch(c *plugin.Context) {}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
