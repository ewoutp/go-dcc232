package dcc232

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func f(x string) string {
	return strings.Replace(x, " ", "", -1)
}

func TestIdlePacket(t *testing.T) {
	p := IdlePacket()
	assert.Equal(t, f("111111111111111 0 11111111 0 00000000 0 11111111 1"), p.String())
}
func TestSpeedAndDirection(t *testing.T) {
	// SpeedSteps: 14 - Direction: Forward
	assert.Equal(t, f("111111111111111 0 00000011 0 01100000 0 01100011 1"), SpeedAndDirection(3, 0, true, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01100011 0 01100010 1"), SpeedAndDirection(1, 2, true, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01101111 0 01101110 1"), SpeedAndDirection(1, 14, true, SpeedSteps14).String())
	// SpeedSteps: 14 - Direction: Reverse
	assert.Equal(t, f("111111111111111 0 00000011 0 01000000 0 01000011 1"), SpeedAndDirection(3, 0, false, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01000011 0 01000010 1"), SpeedAndDirection(1, 2, false, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01001111 0 01001110 1"), SpeedAndDirection(1, 14, false, SpeedSteps14).String())
	// SpeedSteps: 28 - Direction: Forward
	assert.Equal(t, f("111111111111111 0 00000011 0 01100000 0 01100011 1"), SpeedAndDirection(3, 0, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01110010 0 01110011 1"), SpeedAndDirection(1, 2, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01111000 0 01111001 1"), SpeedAndDirection(1, 14, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01101101 0 01101100 1"), SpeedAndDirection(1, 23, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01111111 0 01111110 1"), SpeedAndDirection(1, 28, true, SpeedSteps28).String())
	// SpeedSteps: 28 - Direction: Reverse
	assert.Equal(t, f("111111111111111 0 00000011 0 01000000 0 01000011 1"), SpeedAndDirection(3, 0, false, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01010010 0 01010011 1"), SpeedAndDirection(1, 2, false, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01011000 0 01011001 1"), SpeedAndDirection(1, 14, false, SpeedSteps28).String())
	// 14-bits address, SpeedSteps: 28 - Direction: Reverse
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 01011000 0 01000000 1"), SpeedAndDirection(216, 14, false, SpeedSteps28).String())
}

func TestFunctionGroupOne(t *testing.T) {
	// Address=216, FL=1, F1=0, F2=0, F3=0, F4=0
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 10010000 0 10001000 1"), FunctionGroupOne(216, true, false, false, false, false).String())

}
