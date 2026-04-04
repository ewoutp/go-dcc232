package dcc232

import (
	"fmt"
)

// What we're encoding for
type EncodingTarget byte

const (
	// Target a serial port with an exact match of 1/2 DCC bit equals 1 serial bit.
	EncodingTargetSerialExact EncodingTarget = iota
	// Target a serial port with a high frequency, using multiple serial bytes for a single DCC.
	EncodingTargetSerialHighFreq EncodingTarget = iota
	// Target a SPI port with a high frequency, using multiple serial bytes for a single DCC.
	EncodingTargetSPI EncodingTarget = iota
)

// Baud rate used for this target
func (t EncodingTarget) BaudRate() int {
	switch t {
	case EncodingTargetSerialExact:
		return 17241
	case EncodingTargetSerialHighFreq:
		return 500000
	case EncodingTargetSPI:
		return 500000
	default:
		return 0
	}
}

// Maximum length of encoded packets for this target (in bytes)
func (t EncodingTarget) MaxEncodedLength() int {
	switch t {
	case EncodingTargetSerialExact:
		return ((maxPacketBitLength + 1) / 2)
	case EncodingTargetSerialHighFreq:
		return maxPacketBitLength * 12
	case EncodingTargetSPI:
		return ((maxPacketBitLength * spiBitsPerDCC0 * 2) / 8) + 1
	default:
		return 0
	}
}

const (
	// Maximum length of encoded data, independent of encoding target
	MaxEncodedLength = ((maxPacketBitLength * spiBitsPerDCC0 * 2) / 8) + 1
)

const (
	startPos = 0
	stopPos  = 9
)

const (
	// Bits required at 500kHz (2us per bit)
	spiBitsPerDCC1 = 29
	spiBitsPerDCC0 = 58
)

// EncodePacket takes a DCC packet (as bit stream) and encodes it
// into a byte stream that, when send through a hardware port produces
// the correct DCC bit stream.
func EncodePacket(packet Packet, target EncodingTarget, serialBytes []byte) []byte {
	switch target {
	case EncodingTargetSerialExact:
		return encodePacketSerialExact(packet, serialBytes)
	case EncodingTargetSerialHighFreq:
		return encodePacketSerialHighFreq(packet, serialBytes)
	case EncodingTargetSPI:
		return encodePacketSPI(packet, serialBytes)
	default:
		panic(fmt.Sprintf("Unknown target %v", target))
	}
}

// encodePacketSerialExact takes a DCC packet (as bit stream) and encodes it
// into a byte stream that, when send through a serial port produces
// the correct DCC bit stream using exact encoding.
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
func encodePacketSerialExact(packet Packet, serialBytes []byte) []byte {
	var position int
	var currentByte RS232Byte
	var stretch [maxPacketBitLength]byte
	maxStretch := byte(4)
	packetOffset := 0

	maxEncodedLength := (packet.BitLength() + 1) / 2
	if cap(serialBytes) < maxEncodedLength {
		serialBytes = make([]byte, 0, maxEncodedLength)
	} else {
		serialBytes = serialBytes[:0]
	}

	stretchLast0AndRestart := func(startOffset int) {
		offset := startOffset
		for offset >= 0 {
			if !packet.dccBits[offset] {
				// We found a '0'
				if stretch[offset] < maxStretch {
					// We found a non-stretched '0'
					stretch[offset] = stretch[offset] + 2
					for i := offset + 1; i < packet.BitLength(); i++ {
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

	for packetOffset < packet.BitLength() {
		// Check current byte overflow
		if position == stopPos+1 {
			serialBytes = append(serialBytes, byte(currentByte))
			currentByte = 0
			position = startPos
		}

		value := packet.dccBits[packetOffset]
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
		if packetOffset == packet.BitLength() {
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

// encodePacketSerialHighFreq takes a DCC packet (as bit stream) and encodes it
// into a byte stream that, when send through a serial port produces
// the correct DCC bit stream.
//
// The dataformat of RS232 is as follows:
// Startbits: 1 (always low)
// Databits:  8 LSB first
// Stopbits:  1 (always high)
//
// Speed:     500.000 baud
//
// DCC 1 bit: 01    short low, short high
// encoded as 0x00 0x00 0x00 0xFF 0xFF 0xFF
// DCC 0 bit: 0011  long low, long high. high may be longer
// encoded as 0x00 0x00 0x00 0x00 0x00 0x00 0xFF 0xFF 0xFF 0xFF 0xFF 0xFF
//
// Inputs:
// - packet : DCC packet to encode
// - serialBytes : RS232 output bytes buffer. Can be nil
func encodePacketSerialHighFreq(packet Packet, serialBytes []byte) []byte {
	maxEncodedLength := packet.BitLength() * 12
	if cap(serialBytes) < maxEncodedLength {
		serialBytes = make([]byte, 0, maxEncodedLength)
	} else {
		serialBytes = serialBytes[:0]
	}

	for _, b := range packet.dccBits {
		if b {
			// DCC 1
			serialBytes = append(serialBytes, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF)
		} else {
			// DCC 0
			serialBytes = append(serialBytes, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF)
		}
	}
	return serialBytes
}

// encodePacketSPI takes a DCC packet (as bit stream) and encodes it
// into a byte stream that, when send through an SPI port produces
// the correct DCC bit stream.
//
// Speed:     500.000 baud
//
// DCC 1 bit: 01    short low, short high
// DCC 0 bit: 0011  long low, long high. high may be longer
func encodePacketSPI(packet Packet, serialBytes []byte) []byte {
	maxEncodedLength := ((packet.BitLength() * spiBitsPerDCC0 * 2) / 8) + 1
	if cap(serialBytes) < maxEncodedLength {
		serialBytes = make([]byte, maxEncodedLength)
	} else {
		serialBytes = serialBytes[:maxEncodedLength]
		for idx := range serialBytes {
			serialBytes[idx] = 0
		}
	}
	offset := 0

	setBit := func(high bool) {
		if high {
			byteIdx := offset / 8
			bitIdx := 7 - (offset % 8)
			serialBytes[byteIdx] |= (1 << bitIdx)
		}
		offset++
	}

	for _, b := range packet.dccBits {
		count := spiBitsPerDCC0
		if b {
			count = spiBitsPerDCC1
		}

		// High half-cycle
		for i := 0; i < count; i++ {
			setBit(true)
		}
		// Low half-cycle
		for i := 0; i < count; i++ {
			setBit(false)
		}
	}
	return serialBytes[:(offset+7)/8]
}
