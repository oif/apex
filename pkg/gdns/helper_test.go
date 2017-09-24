package gdns

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHTTPRequest(t *testing.T) {
	HTTPClient = nil
	_, _, err := httpGet(http.MethodGet, "https://github.com", nil)
	assert.Error(t, err, "Error due to did not initialize http client")

	HTTPClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	_, statusCode, err := httpGet(http.MethodGet, "https://github.com", nil)
	assert.Equalf(t, 200, statusCode, "Except status code 200, but %d get", statusCode)
	assert.NoError(t, err, "Error due to did not initialize http client")
}
