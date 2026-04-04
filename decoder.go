package dcc232

import "strings"

// DecodePacket takes a serial byte stream and decodes it into a DCC packet.
//
// The dataformat of RS232 is as follows:
// Startbits: 1 (always low)
// Databits:  8 LSB first
// Stopbits:  1 (always high)
//
// Speed:     19200 baud
//
// DCC 1 bit: 01    short low, short high
// DCC 0 bit: 0011  long low, long high. high may be longer
//
// Inputs:
// - packet : DCC packet to encode
// - serialBytes : RS232 output bytes buffer. Can be nil
func DecodePacket(serialBytes []byte) (string, Packet) {
	result := strings.Builder{}
	for _, sb := range serialBytes {
		result.WriteString("L")
		for i := 0; i < 8; i++ {
			if sb&1 == 0 {
				result.WriteString("L")
			} else {
				result.WriteString("H")
			}
			sb = sb >> 1
		}
		result.WriteString("H")
	}
	fullBitStream := result.String()
	var packet Packet
	remainingBitStream := fullBitStream
	for {
		if remainingBitStream == "" {
			return fullBitStream, packet
		}
		if strings.HasPrefix(remainingBitStream, "LH") {
			packet.dccBits = append(packet.dccBits, true)
			remainingBitStream = remainingBitStream[2:]
			continue
		}
		if strings.HasPrefix(remainingBitStream, "LLHH") {
			packet.dccBits = append(packet.dccBits, false)
			remainingBitStream = remainingBitStream[4:]
			continue
		}
		if strings.HasPrefix(remainingBitStream, "LLLHHH") {
			packet.dccBits = append(packet.dccBits, false)
			remainingBitStream = remainingBitStream[6:]
			continue
		}
		if strings.HasPrefix(remainingBitStream, "LLLLHHHH") {
			packet.dccBits = append(packet.dccBits, false)
			remainingBitStream = remainingBitStream[8:]
			continue
		}
		panic("invalid bitstream: " + remainingBitStream)
	}
}
