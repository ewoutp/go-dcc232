package dcc232

// Packet is a DCC bit stream
type Packet []bool

const (
	preambleCount = 15
)

// encodeByte encoded a given byte into the first 8 positions
// from given offset of the packet.
func (p Packet) encodeByte(offset int, value byte) {
	for i := 7; i >= 0; i-- {
		p[offset+i] = (value & 0x01) == 1
		value = value >> 1
	}
}

// IdlePacket creates an Idle packet
func IdlePacket() Packet {
	p := make(Packet, preambleCount+1+(3*9))
	offset := 0
	for i := 0; i < preambleCount; i++ {
		p[offset] = true
		offset++
	}
	p.encodeByte(offset, 0xFF)
	offset += 9
	p.encodeByte(offset, 0x00)
	offset += 9
	p.encodeByte(offset, 0xFF)
	offset += 8
	p[offset] = true
	return p
}
