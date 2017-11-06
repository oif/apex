package statistics

import (
	"sync"
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

	points    []*client.Point
	writeLock sync.Mutex
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
	if err != nil {
		return err
	}
	// cron job
	go func() {
		for {
			time.Sleep(10 * time.Second)
			p.writeResponse()
		}
	}()
	return nil
}

// Warmup implements plugin
func (p *Plugin) Warmup(c *plugin.Context) {
	c.Set("statistics_plugin:startTime", makeNanoTimestamp())
}

// AfterResponse implements plugin
func (p *Plugin) AfterResponse(c *plugin.Context, err error) {
	responseTime := makeNanoTimestamp() - c.GetInt64("statistics_plugin:startTime")
	c.Logger().WithFields(logrus.Fields{
		"response_time": responseTime,
	}).Info("Response time usage statistics")
	// write influxdb
	go func(responseTime int64) {
		if len(c.Msg.Question) < 1 {
			return
		}
		err := p.pushResponsePoint(c.Msg.Question[0].Qtype, c.Msg.Question[0].Name, responseTime, !c.HasError())
		if err != nil {
			logrus.Error(err)
		}
	}(responseTime)
}

// Patch the dns pakcage
func (p *Plugin) Patch(c *plugin.Context) {}

func makeNanoTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (p *Plugin) pushResponsePoint(qtype uint16, name string, responseTime int64, isSuccess bool) (err error) {
	pt, err := client.NewPoint(
		"dns_query",
		map[string]string{
			"qtype": string(qtype),
			"name":  name,
		},
		map[string]interface{}{
			"response_time": responseTime,
			"success":       isSuccess,
		},
		time.Now(),
	)
	if err != nil {
		return
	}

	p.writeLock.Lock()
	p.points = append(p.points, pt)
	p.writeLock.Unlock()
	return
}

func (p *Plugin) writeResponse() {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "apex",
	})
	if err != nil {
		return
	}

	p.writeLock.Lock()
	for _, p := range p.points {
		bp.AddPoint(p)
	}
	p.points = make([]*client.Point, 0)
	p.writeLock.Unlock()

	err = influxdbClient.Write(bp) // ignore error
	if err != nil {
		logrus.Error(err)
	}
}
