package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApexd(t *testing.T) {
	// Current has nothing should test
	assert.NotPanics(t, func() {
		main()
	}, "No panic")
}
