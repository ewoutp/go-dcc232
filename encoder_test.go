package dcc232

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func bits(encoded []byte) string {
	result := ""
	for idx, b := range encoded {
		if idx > 0 {
			result += " "
		}
		result = result + RS232Byte(b).String()
	}
	return result
}

func TestEncodePacket(t *testing.T) {
	p := IdlePacket()
	sb := EncodePacket(p)
	assert.Equal(t, "5555555655c6c6c6c6665555", hex.EncodeToString(sb))
	assert.Equal(t, "0101010101 0101010101 0101010101 0011010101 0101010101 0011000111 0011000111 0011000111 0011000111 0011001101 0101010101 0101010101", bits(sb))
	assert.Equal(t, "43 / 12", fmt.Sprintf("%d / %d", len(p), len(sb)))
}

func TestEncodeSpeed(t *testing.T) {
	p := SpeedAndDirection(1, 2, true, SpeedSteps128)
	sb := EncodePacket(p)
	assert.Equal(t, "0101010101 0101010101 0101010101 0011000111 0011000111 0011000111 0011001101 0011000111 0011010101 0101010011 0100110011 0011000111 0011010011 0000111101 0011010101 0100110011 0101010101", bits(sb))
	assert.Equal(t, "52 / 17", fmt.Sprintf("%d / %d", len(p), len(sb)))
}

func TestEncodeAllSpeeds(t *testing.T) {
	for addr := 1; addr < 200; addr++ {
		for speed := byte(0); speed < 128; speed++ {
			p := SpeedAndDirection(addr, speed, true, SpeedSteps128)
			func() {
				defer func() {
					if err := recover(); err != nil {
						t.Fatalf("EncodePacket failed for addr=%d, speed=%d packet=%s (err=%s)", addr, speed, p.String(), err)
					}

				}()
				EncodePacket(p)
			}()
		}
	}
}

func TestEncodeFunctionGroupOne(t *testing.T) {
	p := FunctionGroupOne(216, true, false, false, false, false)
	sb := EncodePacket(p)
	assert.Equal(t, f("111111111111111 0 11000000 0 11011000 0 10010000 0 10001000 1"), p.String())
	assert.Equal(t, "0101010101 0101010101 0101010101 0001110101 0011000111 0011000111 0011000111 0001110101 0001110101 0011000111 0011001101 0011001101 0011000111 0011000111 0011010011 0011001101 0011000111 0011010101", bits(sb))
}
