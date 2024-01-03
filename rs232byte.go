package dcc232

import "fmt"

// RS232Byte is a helper for building up RS232 data streams.
type RS232Byte uint8

const (
	startBit = false
	stopBit  = true
)

// Get the bit at given position.
// Positions:
// 0: start bit, always false
// 1..8: data bites
// 9: stop bit, always true
// This function panics in index is out of bounds.
func (b RS232Byte) Get(index int) bool {
	if index < 0 || index > 9 {
		panic(fmt.Sprintf("index (%d) out of bounds", index))
	}
	switch index {
	case 0:
		return startBit
	case 9:
		return stopBit
	default:
		index--
		return (byte(b) & (1 << uint(index))) != 0
	}
}

// Set the bit at given position.
// Positions:
// 0: start bit, must be false
// 1..8: data bites
// 9: stop bit, must be true
// This function panics in index is out of bounds or if a required value
// is incorrect.
func (b *RS232Byte) Set(index int, value bool) {
	if index < 0 || index > 9 {
		panic(fmt.Sprintf("index (%d) out of bounds", index))
	}
	switch index {
	case 0:
		if value != startBit {
			panic("startBit expected")
		}
	case 9:
		if value != stopBit {
			panic("stopBit expected")
		}
	default:
		index--
		bit := 1 << uint(index)
		if value {
			*b = *b | RS232Byte(bit)
		} else {
			*b = *b & ^RS232Byte(bit)
		}
	}
}

// Convert the given byte into a bit string (including start & stop bits)
func (b RS232Byte) String() string {
	result := make([]byte, 10)
	for i := 0; i < 10; i++ {
		if b.Get(i) {
			result[i] = '1'
		} else {
			result[i] = '0'
		}
	}
	return string(result)
}
