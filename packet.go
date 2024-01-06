package dcc232

import "fmt"

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

// Packet length
const (
	// Number of preamble bits
	preambleCount = 15

	// Maximum number of bytes takes by an address
	MaxAddressBytes = 2

	// Maximum number of bytes taken by the data portion of a packet
	MaxDataBytes = 3

	MaxPacketLength = preambleCount + 1 + (MaxAddressBytes * 9) + (MaxDataBytes * 9) + 9
)

// SpeedSteps of a loc decoder
type SpeedSteps int

const (
	SpeedSteps128 SpeedSteps = 128
	SpeedSteps28  SpeedSteps = 28
	SpeedSteps14  SpeedSteps = 14
)

var (
	speed28Steps = []byte{
		0b00000000, // Stop
		0b00000010, // Step 1
		0b00010010, // Step 2
		0b00000011, // Step 3
		0b00010011, // Step 4
		0b00000100, // Step 5
		0b00010100, // Step 6
		0b00000101, // Step 7
		0b00010101, // Step 8
		0b00000110, // Step 9
		0b00010110, // Step 10
		0b00000111, // Step 11
		0b00010111, // Step 12
		0b00001000, // Step 13
		0b00011000, // Step 14
		0b00001001, // Step 15
		0b00011001, // Step 16
		0b00001010, // Step 17
		0b00011010, // Step 18
		0b00001011, // Step 19
		0b00011011, // Step 20
		0b00001100, // Step 21
		0b00011100, // Step 22
		0b00001101, // Step 23
		0b00011101, // Step 24
		0b00001110, // Step 25
		0b00011110, // Step 26
		0b00001111, // Step 27
		0b00011111, // Step 28
	}
)

// ensureMaxLength checks the capacity of the given packet.
// If large enough for an maximum size packet, the given
// packet is set to maximum length and return.
// Otherwise a new packet is allocated and returned.
func (p Packet) ensureMaxLength() Packet {
	if cap(p) >= MaxPacketLength {
		// Set to max length
		return p[:]
	}
	return make(Packet, MaxPacketLength)
}

// encodePreamble writes the preamble into the head of the given packet
// including a trailing '0'.
// Returns: offset for first address byte.
func (p Packet) encodePreamble() int {
	offset := 0
	// preamble '111111111111111'
	for i := 0; i < preambleCount; i++ {
		p[offset] = true
		offset++
	}
	// '0'
	offset++
	p[offset] = false
	return offset
}

// encodeAddress writes the given address into the given packet
// at the given offset, including a trailing '0'.
// Returns: offset after address+'0', updated error code
func (p Packet) encodeAddress(address, offset int, error byte) (int, byte) {
	addressBytes := 1
	if address > 127 {
		addressBytes = 2
	}

	// First address byte
	if addressBytes == 1 {
		// Single
		value := byte(address & 0x7f)
		error = p.encodeByte(offset, value, error)
		offset += 8
		p[offset] = false
		offset++
	} else {
		// 2 address bytes
		value1 := byte(0xc0 | ((address >> 8) & 0x3f))
		error = p.encodeByte(offset, value1, error)
		offset += 8
		p[offset] = false
		offset++
		value2 := byte(address & 0xff)
		error = p.encodeByte(offset, value2, error)
		offset += 8
		p[offset] = false
		offset++
	}
	return offset, error
}

// encodeByte encoded a given byte into the first 8 positions
// from given offset of the packet.
// Returns: updated error code
func (p Packet) encodeByte(offset int, value, error byte) byte {
	result := error ^ value
	for i := 7; i >= 0; i-- {
		p[offset+i] = (value & 0x01) == 1
		value = value >> 1
	}
	return result
}

// IdlePacket creates an Idle packet
func (p Packet) IdlePacket() Packet {
	// Use packet of maximum length
	p = p.ensureMaxLength()

	// Preamble
	offset := p.encodePreamble()

	// '11111111 0'
	error := p.encodeByte(offset, 0xFF, 0)
	offset += 9
	// '00000000 0'
	error = p.encodeByte(offset, 0x00, error)
	offset += 9
	// '11111111 1'
	p.encodeByte(offset, error, 0)
	offset += 8
	p[offset] = true
	offset++
	return p[:offset]
}

// SpeedAndDirection creates a standard speed & direction packet
func (p Packet) SpeedAndDirection(address int, speed byte, direction bool, speedSteps SpeedSteps) Packet {
	// Use packet of maximum length
	p = p.ensureMaxLength()

	// Preamble
	offset := p.encodePreamble()
	error := byte(0)

	// Address
	offset, error = p.encodeAddress(address, offset, error)

	if speedSteps != SpeedSteps128 {
		// 14 or 28 speed steps
		// data byte
		data := byte(0b01000000)
		if direction {
			data |= 0b00100000
		}
		if speedSteps == SpeedSteps14 {
			if speed > 0 {
				// Skip E-Stop
				speed++
			}
			speed &= 0x0f
		} else {
			if speed >= 0 && speed <= 28 {
				speed = speed28Steps[speed]
			} else {
				speed = 0
			}
		}
		data |= speed

		error = p.encodeByte(offset, data, error)
		offset += 8
		p[offset] = false
		offset++
	} else {
		// 128 speed steps
		// data byte
		data1 := byte(0x3f)
		data2 := byte(speed & 0x7f)
		if direction {
			data2 |= 0x80
		}

		error = p.encodeByte(offset, data1, error)
		offset += 8
		p[offset] = false
		offset++
		error = p.encodeByte(offset, data2, error)
		offset += 8
		p[offset] = false
		offset++
	}
	p.encodeByte(offset, error, 0)
	offset += 8
	p[offset] = true
	offset++
	return p[:offset]
}

// FunctionGroupOne creates a packet to control F0, F1-F4
func (p Packet) FunctionGroupOne(address int, fl, f1, f2, f3, f4 bool) Packet {
	// Use packet of maximum length
	p = p.ensureMaxLength()

	// Preamble
	offset := p.encodePreamble()
	error := byte(0)

	// Address
	offset, error = p.encodeAddress(address, offset, error)

	// Data
	data := byte(0x80)
	if f1 {
		data |= 0x01
	}
	if f2 {
		data |= 0x02
	}
	if f3 {
		data |= 0x04
	}
	if f4 {
		data |= 0x08
	}
	if fl {
		data |= 0x10
	}
	error = p.encodeByte(offset, data, error)
	offset += 8
	p[offset] = false
	offset++

	p.encodeByte(offset, error, 0)
	offset += 8
	p[offset] = true
	offset++
	return p[:offset]
}

// FunctionGroupTwo creates a packet to control F5-F8 (firstIndex=5) or F9-F12 (firstIndex=9)
func (p Packet) FunctionGroupTwo(address int, firstIndex byte, fa, fb, fc, fd bool) Packet {
	data := byte(0b10100000)
	switch firstIndex {
	case 5:
		data |= 0b00010000
	case 9:
		// Do nothing
	default:
		panic(fmt.Errorf("invalid firstIndex %d", firstIndex))
	}
	// Use packet of maximum length
	p = p.ensureMaxLength()

	// Preamble
	offset := p.encodePreamble()
	error := byte(0)

	// Address
	offset, error = p.encodeAddress(address, offset, error)

	// Data
	if fa {
		data |= 0x01
	}
	if fb {
		data |= 0x02
	}
	if fc {
		data |= 0x04
	}
	if fd {
		data |= 0x08
	}
	error = p.encodeByte(offset, data, error)
	offset += 8
	p[offset] = false
	offset++

	p.encodeByte(offset, error, 0)
	offset += 8
	p[offset] = true
	offset++
	return p[:offset]
}
