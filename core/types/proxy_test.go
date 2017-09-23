package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyCheck(t *testing.T) {
	p := &Proxy{}
	assert.Error(t, p.Check(), "Error due to empty protocol")

	p.Protocol = "http"
	assert.Error(t, p.Check(), "Error due to empty address")

	p.Address = "127.0.0.1:1080"
	assert.NoError(t, p.Check(), "No error")
}
