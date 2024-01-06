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
	p := make(Packet, MaxPacketLength)
	p = p.IdlePacket()
	assert.Equal(t, f("111111111111111 0 11111111 0 00000000 0 11111111 1"), p.String())
}
func TestSpeedAndDirection(t *testing.T) {
	p := make(Packet, MaxPacketLength)

	// SpeedSteps: 14 - Direction: Forward
	assert.Equal(t, f("111111111111111 0 00000011 0 01100000 0 01100011 1"), p.SpeedAndDirection(3, 0, true, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01100011 0 01100010 1"), p.SpeedAndDirection(1, 2, true, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01101111 0 01101110 1"), p.SpeedAndDirection(1, 14, true, SpeedSteps14).String())
	// SpeedSteps: 14 - Direction: Reverse
	assert.Equal(t, f("111111111111111 0 00000011 0 01000000 0 01000011 1"), p.SpeedAndDirection(3, 0, false, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01000011 0 01000010 1"), p.SpeedAndDirection(1, 2, false, SpeedSteps14).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01001111 0 01001110 1"), p.SpeedAndDirection(1, 14, false, SpeedSteps14).String())
	// SpeedSteps: 28 - Direction: Forward
	assert.Equal(t, f("111111111111111 0 00000011 0 01100000 0 01100011 1"), p.SpeedAndDirection(3, 0, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01110010 0 01110011 1"), p.SpeedAndDirection(1, 2, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01111000 0 01111001 1"), p.SpeedAndDirection(1, 14, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01101101 0 01101100 1"), p.SpeedAndDirection(1, 23, true, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01111111 0 01111110 1"), p.SpeedAndDirection(1, 28, true, SpeedSteps28).String())
	// SpeedSteps: 28 - Direction: Reverse
	assert.Equal(t, f("111111111111111 0 00000011 0 01000000 0 01000011 1"), p.SpeedAndDirection(3, 0, false, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01010010 0 01010011 1"), p.SpeedAndDirection(1, 2, false, SpeedSteps28).String())
	assert.Equal(t, f("111111111111111 0 00000001 0 01011000 0 01011001 1"), p.SpeedAndDirection(1, 14, false, SpeedSteps28).String())
	// 14-bits address, SpeedSteps: 28 - Direction: Reverse
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 01011000 0 01000000 1"), p.SpeedAndDirection(216, 14, false, SpeedSteps28).String())
}

func TestFunctionGroupOne(t *testing.T) {
	p := make(Packet, MaxPacketLength)

	// Address=216, FL=1, F1=0, F2=0, F3=0, F4=0
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 10010000 0 10001000 1"), p.FunctionGroupOne(216, true, false, false, false, false).String())
	// Address=83, FL=0, F1=1, F2=1, F3=1, F4=0
	assert.Equal(t, f("111111111111111 0 01010011 0 10000111 0 11010100 1"), p.FunctionGroupOne(83, false, true, true, true, false).String())
	// Address=83, FL=0, F1=1, F2=0, F3=0, F4=1
	assert.Equal(t, f("111111111111111 0 01010011 0 10001001 0 11011010 1"), p.FunctionGroupOne(83, false, true, false, false, true).String())
}

func TestFunctionGroupTwo(t *testing.T) {
	var p Packet

	// Address=216, F5=0, F6=0, F7=0, F8=0
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 10110000 0 10101000 1"), p.FunctionGroupTwo(216, 5, false, false, false, false).String())
	// Address=83, F5=1, F6=1, F7=1, F8=0
	assert.Equal(t, f("111111111111111 0 01010011 0 10110111 0 11100100 1"), p.FunctionGroupTwo(83, 5, true, true, true, false).String())
	// Address=83, F5=1, F6=0, F7=0, F8=1
	assert.Equal(t, f("111111111111111 0 01010011 0 10111001 0 11101010 1"), p.FunctionGroupTwo(83, 5, true, false, false, true).String())

	// Address=216, F9=0, F10=0, F11=0, F12=0
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 10100000 0 10111000 1"), p.FunctionGroupTwo(216, 9, false, false, false, false).String())
	// Address=83, F9=1, F10=1, F11=1, F12=0
	assert.Equal(t, f("111111111111111 0 01010011 0 10100111 0 11110100 1"), p.FunctionGroupTwo(83, 9, true, true, true, false).String())
	// Address=83, F9=1, F10=0, F11=0, F12=1
	assert.Equal(t, f("111111111111111 0 01010011 0 10101001 0 11111010 1"), p.FunctionGroupTwo(83, 9, true, false, false, true).String())
}
