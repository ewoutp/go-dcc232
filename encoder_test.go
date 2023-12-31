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
	assert.Equal(t, "5555555595e6e6e6e65655f9", hex.EncodeToString(sb))
	assert.Equal(t, "0101010101 0101010101 0101010101 0101010101 0101010011 0011001111 0011001111 0011001111 0011001111 0011010101 0101010101 0100111111", bits(sb))
	assert.Equal(t, "43 / 12", fmt.Sprintf("%d / %d", len(p), len(sb)))
}

func TestEncodeSpeed(t *testing.T) {
	p := SpeedAndDirection(1, 2, true, SpeedSteps128)
	sb := EncodePacket(p)
	assert.Equal(t, "0101010101 0101010101 0101010101 0011001111 0011001111 0011001111 0011001101 0011001111 0011010101 0101010011 0100110011 0011001111 0011010011 0011111101 0011010101 0100110011 0101010101", bits(sb))
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
