package dcc232

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func f(x string) string {
	return strings.Replace(x, " ", "", -1)
}

func TestSpeedAndDirection(t *testing.T) {
	// SpeedSteps14
	assert.Equal(t, f("111111111111111 0 00000011 0 01100000 0 01100011 1"), SpeedAndDirection(3, 0, true, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01100010 0 01100011 1"), SpeedAndDirection(1, 2, true, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01101110 0 01101111 1"), SpeedAndDirection(1, 14, true, SpeedSteps14).String())
}
