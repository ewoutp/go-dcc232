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
	// preamble '111111111111111'
	for i := 0; i < preambleCount; i++ {
		p[offset] = true
		offset++
	}
	// '0'
	offset++
	// '11111111 0'
	p.encodeByte(offset, 0xFF)
	offset += 9
	// '00000000 0'
	p.encodeByte(offset, 0x00)
	offset += 9
	// '11111111 1'
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

// FunctionGroupOne creates a packet to control F0, F1-F4
func FunctionGroupOne(address int, fl, f1, f2, f3, f4 bool) Packet {
	addressBytes := 1
	if address > 127 {
		addressBytes = 2
	}
	dataBytes := 1
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
	bits.encodeByte(offset, data)
	offset += 9
	error ^= data

	bits.encodeByte(offset, error)
	offset += 8
	bits[offset] = true
	return bits
}
