package plugin

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/oif/apex/pkg/gdns"
	plugin "github.com/oif/apex/pkg/plugin/v1"
)

// GoogleDNSPluginName for g.Name
const GoogleDNSPluginName = "Google DNS Plugin"

// GoogleDNS plugin implements pkg/plugin/v1
type GoogleDNS struct{}

// Name return the name of this plugin
func (g *GoogleDNS) Name() string {
	return GoogleDNSPluginName
}

// Initialize Google DNS
func (g *GoogleDNS) Initialize() error {
	proxyAddr, _ := url.Parse("http://127.0.0.1:6152")

	gdns.HTTPClient = &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyAddr),
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return nil
}

// Patch the dns pakcage
func (g *GoogleDNS) Patch(c *plugin.Context) {
	// construct google dns request body
	rr := new(gdns.ResolveRequest)
	// get first question default
	if len(c.Msg.Question) < 1 {
		// abort due to no question here
		c.Abort()
		return
	}
	question := c.Msg.Question[0]
	rr.Name = question.Name
	rr.Type = question.Qtype
	resp, _, err := rr.Request() // ignore status code current
	if err != nil {
		c.AbortWithError(err)
		return
	}
	response, err := gdns.BytesToResolveResponse(resp)
	// json decode error
	if err != nil {
		c.AbortWithError(err)
		return
	}
	// resolve error
	if ok, comment := response.Success(); !ok {
		c.AbortWithError(errors.New(comment))
		return
	}

	for _, ans := range response.Answer {
		// construct every response for dnsPack
		c.Msg.Answer = append(c.Msg.Answer, ans.ToRR())
	}
}
