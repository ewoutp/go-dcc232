package dcc232

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodePacket(t *testing.T) {
	var src Packet
	srcPkt := src.SpeedAndDirection(36, 122, true, SpeedSteps128, false)
	srcEncoded := EncodePacket(srcPkt, EncodingTargetSerialExact, nil)
	assert.Equal(t, "1111111111111111111111 0 00100100 0 00111111 0 11111011 0 11100000 1", srcPkt.String())

	encoded := hex.EncodeToString(srcEncoded)
	serialBytes, err := hex.DecodeString(encoded)
	require.NoError(t, err)
	bitStream, decoded := DecodePacket(serialBytes)
	assert.Equal(t, "1111111111111111111111 0 00100100 0 00111111 0 11111011 0 11100000 1 11", decoded.String(), bitStream)
}
