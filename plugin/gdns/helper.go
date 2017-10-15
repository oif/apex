package gdns

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// GoogleDNSAPI Google DNS over HTTPS
	// ref. https://developers.google.com/speed/public-dns/docs/dns-over-https
	GoogleDNSAPI = "https://dns.google.com/"
	// RequestRetryBackoff default 3
	RequestRetryBackoff = 3
)

var (
	// HTTPClient used to query Google DNS API
	HTTPClient *http.Client
)

func httpGet(method, endpoint string, urlParams map[string]interface{}) (response []byte, statusCode int, err error) {
	if HTTPClient == nil {
		err = errors.New("http client didn't intialized")
		return
	}

	retryTimes := 0
PAYLOAD:
	response, statusCode, err = clientRequest(method, endpoint, urlParams, nil)
	if (statusCode >= 400 || err != nil) && retryTimes < RequestRetryBackoff {
		// retry request
		retryTimes++
		goto PAYLOAD
	}
	return
}

func clientRequest(method, endpoint string, urlParams map[string]interface{}, body io.Reader) (response []byte, statusCode int, err error) {
	if params := paramsFormator(urlParams); len(params) != 0 {
		endpoint += ("?" + params)
	}
	var request *http.Request
	request, err = http.NewRequest(method, endpoint, body)
	if err != nil {
		return
	}

	resp, err := HTTPClient.Do(request)
	if err != nil {
		return
	}
	statusCode = resp.StatusCode
	defer resp.Body.Close()

	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func paramsFormator(params map[string]interface{}) string {
	if params == nil || len(params) == 0 {
		// no params given
		return ""
	}
	var result string
	for key, val := range params {
		result += fmt.Sprintf("%s=%v&", key, val)
	}
	return strings.TrimRight(result, "&")
}
