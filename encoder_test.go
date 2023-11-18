package dcc232

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePacket(t *testing.T) {
	p := IdlePacket()
	sb := EncodePacket(p)
	assert.Equal(t, "5555555595e6e6e6e6565519", hex.EncodeToString(sb))
	assert.Equal(t, "43 / 12", fmt.Sprintf("%d / %d", len(p), len(sb)))
}

func TestEncodeSpeed(t *testing.T) {
	p := SpeedAndDirection(1, 2, true, SpeedSteps128)
	sb := EncodePacket(p)
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
