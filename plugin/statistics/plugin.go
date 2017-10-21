package statistics

import (
	"time"

	plugin "github.com/oif/apex/pkg/plugin/v1"

	"github.com/Sirupsen/logrus"
	client "github.com/influxdata/influxdb/client/v2"
)

var influxdbClient client.Client

// PluginName for g.Name
const PluginName = "Statistics Plugin"

// Plugin implements pkg/plugin/v1
type Plugin struct {
	ConfigFilePath string
}

// Name return the name of this plugin
func (p *Plugin) Name() string {
	return PluginName
}

// Initialize Google DNS Plugin
func (p *Plugin) Initialize() error {
	c := new(Config)
	c.Load(p.ConfigFilePath)
	var err error
	influxdbClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     c.InfluxDB.Addr,
		Username: c.InfluxDB.Username,
		Password: c.InfluxDB.Password,
	})

	return err
}

// Warmup implements plugin
func (p *Plugin) Warmup(c *plugin.Context) {
	c.Set("statistics_plugin:startTime", makeNanoTimestamp())
}

// AfterResponse implements plugin
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {
	if startAt := c.GetInt64("statistics_plugin:startTime"); startAt != 0 {
		responseTime := makeNanoTimestamp() - startAt
		c.Logger().WithFields(logrus.Fields{
			"response_time(ns)": responseTime,
		}).Info("Response time usage statistics")
		// write influxdb
		go func() {
			err := writeResponse(c.Msg.Question[0].Qtype, c.Msg.Question[0].Name, responseTime, !c.HasError())
			if err != nil {
				c.Logger().WithFields(logrus.Fields{
					"err": err,
				}).Error("Write influxdb error")
			}
		}()
	}
}

// Patch the dns pakcage
func (p *Plugin) Patch(c *plugin.Context) {}

func makeNanoTimestamp() int64 {
	return time.Now().UnixNano()
}

func writeResponse(qtype uint16, name string, responseTime int64, isSuccess bool) (err error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "apex",
	})
	if err != nil {
		return
	}

	tags := map[string]string{
		"qtype": string(qtype),
		"name":  name,
	}

	fields := map[string]interface{}{
		"response_time": responseTime,
		"success":       isSuccess,
	}

	pt, err := client.NewPoint(
		"dns_query",
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		return
	}
	bp.AddPoint(pt)

	return influxdbClient.Write(bp)
}
