package dcc232

import "fmt"

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
	stretch := make([]byte, len(packet))
	maxStretch := byte(4)
	packetOffset := 0

	stretchLast0AndRestart := func(startOffset int) {
		offset := startOffset
		for offset >= 0 {
			if !packet[offset] {
				// We found a '0'
				if stretch[offset] < maxStretch {
					// We found a non-stretched '0'
					stretch[offset] = stretch[offset] + 2
					for i := offset + 1; i < len(packet); i++ {
						stretch[i] = 0
					}
					packetOffset = 0
					serialBytes = serialBytes[:0]
					position = 0
					currentByte = 0
					return
				}
			}
			offset--
		}
		// Find last '0'
		panic(fmt.Sprintf("No unstretched '0' bit available starting at %d in %s", startOffset, packet.String()))
	}

	for packetOffset < len(packet) {
		// Check current byte overflow
		if position == stopPos+1 {
			serialBytes = append(serialBytes, byte(currentByte))
			currentByte = 0
			position = startPos
		}

		value := packet[packetOffset]
		stretched := stretch[packetOffset]
		packetOffset++
		if value {
			// "1"
			// We always have room for "1" bit
			currentByte.Set(position+0, false)
			currentByte.Set(position+1, true)
			position += 2
		} else {
			// "0"
			length := int(4 + stretched)
			if position+length <= 10 {
				// We have room for "0" bit

				// Set bits
				currentByte.Set(position+0, false)
				currentByte.Set(position+1, false)
				currentByte.Set(position+2, true)
				currentByte.Set(position+3, true)
				position += 4
				for i := byte(0); i < stretched; i++ {
					currentByte.Set(position, true)
					position++
				}
			} else {
				// Go back to last "0" and make it longer
				stretchLast0AndRestart(packetOffset - 2)
				continue
			}
		}
		// End of packet, then pad current byte if needed
		if packetOffset == len(packet) {
			// End of package
			for position < stopPos+1 {
				// We need to pad the current byte
				if value {
					// Last value was '1', pad with more '1's
					currentByte.Set(position+0, false)
					currentByte.Set(position+1, true)
					position += 2
				} else {
					// Last value was '0', make it longer
					currentByte.Set(position+0, true)
					currentByte.Set(position+1, true)
					position += 2
				}
			}
		}
	}

	// Add last byte
	if position != startPos {
		serialBytes = append(serialBytes, byte(currentByte))
	}

	return serialBytes
}
