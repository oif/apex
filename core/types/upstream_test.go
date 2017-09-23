package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpstreamCheck(t *testing.T) {
	u := &Upstream{}
	assert.Error(t, u.Check(), "Error due to empty address")
	u.Address = "8.8.8.8:53"

	assert.Error(t, u.Check(), "Error due to empty protocol")
	u.Protocol = "udp"

	assert.Error(t, u.Check(), "Error due to empty DNS name")
	u.Name = "Google DNS"

	u.Timeout = -1
	assert.Error(t, u.Check(), "Error due to invalid timeout")

	u.Timeout = 0
	u.Check()
	assert.NotEqualf(t, DefaultUpstreamTimeout, u.Timeout, "Timeout should be %d but %d get", DefaultUpstreamTimeout, u.Timeout)

	assert.NoError(t, u.Check(), "No error")
}
