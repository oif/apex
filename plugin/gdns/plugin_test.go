package gdns

import (
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func initClient() {
	isTravis := os.Getenv("TRAVIS")

	if isTravis == "true" {
		HTTPClient = &http.Client{
			Timeout: 2 * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
	} else {
		proxyAddr, _ := url.Parse("http://127.0.0.1:6152")

		HTTPClient = &http.Client{
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
	}
}

func TestGoogleDNSPlugin(t *testing.T) {
	request := &ResolveRequest{}
	_, _, err := request.Request()
	assert.Error(t, err, "Invalid resolve name")

	request.Name = "google.com"
	request.Type = dns.TypeNone
	request.Request()
	assert.Equal(t, dns.TypeANY, request.Type, "Default rr type ANY")

	initClient()
	request = &ResolveRequest{
		Name:             "example.com",
		Type:             dns.TypeA,
		CheckingDisabled: true,
	}
	resp, statusCode, err := request.Request()
	assert.NoError(t, err, "Should have no any error")
	assert.Equal(t, 200, statusCode, "Response status code should be 200 but %d get", statusCode)
	// end resolve test

	// start test resolve response
	response, err := BytesToResolveResponse(resp)
	assert.NoError(t, err, "JSON unmarshal error")
	success, comment := response.Success()
	assert.Equal(t, true, success, "Resolve fail", comment)
	assert.Equal(t, dns.TypeA, response.Question[0].Type, "Question type")
}
