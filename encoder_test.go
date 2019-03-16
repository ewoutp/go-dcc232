package dcc232

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodePacket(t *testing.T) {
	p := IdlePacket()
	sb := EncodePacket(p)
	assert.Equal(t, "5555555595e6e6e6e6565519", hex.EncodeToString(sb))
}
