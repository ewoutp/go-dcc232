package dcc232

const (
	startPos = 0
	stopPos  = 9
)

// EncodePacket takes a DCC packet (as bit stream) and encodes it
// into a byte stream that, when send through a serial port produces
// the correct DCC bit stream.
//
// The dataformat of RS232 is as follows:
// Startbits: 1 (always low)
// Stopbits:  1 (always high)
// Databits:  8
// Speed:     19200 baud
//
// DCC 1 bit: 01    short low, short high
// DCC 0 bit: 0011  long low, long high. high may be longer
func EncodePacket(packet Packet) []byte {

	var serialBytes []byte
	var position int
	var currentByte RS232Byte
	packetOffset := 0
	last0PacketOffset := -1
	last0Stretch := 0
	last0BytesLength := 0
	last0Position := 0

	restartAtLast0 := func(stretch int) {
		if last0PacketOffset < 0 {
			panic("No earlier ')' bit available")
		}
		// Set packet offset
		packetOffset = last0PacketOffset + 1
		// Reset current byte
		if len(serialBytes) == last0BytesLength {
			// Still in same current byte
		} else {
			// Restore current byte and trim serialBytes
			currentByte = RS232Byte(serialBytes[last0BytesLength])
			serialBytes = serialBytes[:last0BytesLength]
		}
		// Reset position
		position = last0Position
		// Set stretched bits
		currentByte.Set(position+0, false)
		currentByte.Set(position+1, false)
		currentByte.Set(position+2, true)
		currentByte.Set(position+3, true)
		position += 4
		last0Stretch += stretch
		for i := 0; i < last0Stretch; i++ {
			currentByte.Set(position, true)
			position++
		}
	}

	for packetOffset < len(packet) {
		// Check current byte overflow
		if position == stopPos+1 {
			serialBytes = append(serialBytes, byte(currentByte))
			currentByte = 0
			position = startPos
		}

		value := packet[packetOffset]
		packetOffset++
		if value {
			// "1"
			if position == stopPos {
				// Go back to last "0"
				restartAtLast0(1)
			} else {
				// We have room for "1" bit
				currentByte.Set(position+0, false)
				currentByte.Set(position+1, true)
				position += 2
			}
		} else {
			// "0"
			if position <= 6 {
				// We have room for "0" bit

				// Record position
				last0PacketOffset = packetOffset - 1
				last0Stretch = 0
				last0BytesLength = len(serialBytes)
				last0Position = position

				// Set bits
				currentByte.Set(position+0, false)
				currentByte.Set(position+1, false)
				currentByte.Set(position+2, true)
				currentByte.Set(position+3, true)
				position += 4
			} else {
				// Go back to last "0"
				restartAtLast0(10 - position)
			}
		}
	}

	// Add last byte
	if position != startPos {
		serialBytes = append(serialBytes, byte(currentByte))
	}

	return serialBytes
}
