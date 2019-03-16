package dcc232

// Packet is a DCC bit stream
type Packet []bool

// String converts the packet to a binary string
func (p Packet) String() string {
	s := make([]byte, len(p))
	for i, v := range p {
		if v {
			s[i] = '1'
		} else {
			s[i] = '0'
		}
	}
	return string(s)
}

const (
	preambleCount = 15
)

// SpeedSteps of a loc decoder
type SpeedSteps int

const (
	SpeedSteps128 SpeedSteps = 128
	SpeedSteps28  SpeedSteps = 28
	SpeedSteps14  SpeedSteps = 14
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

// SpeedAndDirection creates a standard speed & direction packet
func SpeedAndDirection(address int, speed byte, direction bool, speedSteps SpeedSteps) Packet {
	addressBytes := 1
	if address > 127 {
		addressBytes = 2
	}
	dataBytes := 1
	if speedSteps == SpeedSteps128 {
		dataBytes = 2
	}
	bits := make(Packet, preambleCount+1+((addressBytes+dataBytes+1)*9))
	offset := preambleCount + 1
	error := byte(0)
	for i := 0; i < preambleCount; i++ {
		bits[i] = true
	}

	// First address byte
	if addressBytes == 1 {
		// Single
		value := byte(address & 0x7f)
		bits.encodeByte(offset, value)
		offset += 9
		error ^= value
	} else {
		// 2 address bytes
		value1 := byte(0xc0 | ((address >> 8) & 0x3f))
		bits.encodeByte(offset, value1)
		offset += 9
		value2 := byte(address & 0xff)
		bits.encodeByte(offset, value2)
		offset += 9
		error ^= value1
		error ^= value2
	}

	if speedSteps != SpeedSteps128 {
		// 14 or 28 speed steps
		// data byte
		data := byte(0x40)
		if direction {
			data |= 0x20
		}
		if speedSteps == SpeedSteps14 {
			speed &= 0x0f
		} else {
			speed &= 0x1f
		}
		data |= speed

		bits.encodeByte(offset, data)
		offset += 9
		error ^= data
	} else {
		// 128 speed steps
		// data byte
		data1 := byte(0x3f)
		data2 := byte(speed & 0x7f)
		if direction {
			data2 |= 0x80
		}

		bits.encodeByte(offset, data1)
		offset += 9
		bits.encodeByte(offset, data2)
		offset += 9
		error ^= data1
		error ^= data2
	}
	bits.encodeByte(offset, error)
	offset += 8
	bits[offset] = true
	return bits
}
