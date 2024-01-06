package dcc232

import "fmt"

const (
	startPos = 0
	stopPos  = 9

	// Maximum length of an encoded packet
	// 1 RS232 byte contains at least 2 DCC bits.
	MaxEncodedLength = (MaxPacketLength / 2)
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
//
// Inputs:
// - packet : DCC packet to encode
// - serialBytes : RS232 output bytes buffer. Can be nil
func EncodePacket(packet Packet, serialBytes []byte) []byte {
	var position int
	var currentByte RS232Byte
	var stretch [MaxPacketLength]byte
	maxStretch := byte(4)
	packetOffset := 0

	maxEncodedLength := (len(packet) + 1) / 2
	if cap(serialBytes) < maxEncodedLength {
		serialBytes = make([]byte, 0, maxEncodedLength)
	} else {
		serialBytes = serialBytes[:0]
	}

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
				for i := 0; i < length; i++ {
					currentByte.Set(position, !(i < length/2))
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
