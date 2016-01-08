// package utf8 decode a Reader for UTF-8 encoded byte streams
package utf8

import (
	"fmt"
	"io"
)

// In November 2003, UTF-8 was restricted by RFC 3629 to end at U+10FFFF,
// in order to match the constraints of the UTF-16 character encoding.
// This removed all 5- and 6-byte sequences, and 983040 4-byte sequences.
// https://en.wikipedia.org/wiki/UTF-8#Description
const MaxBytes = 4

// Reads the next Unicode code point from a reader.
// 
// Skips to the next leading byte if the first byte isn't a leading byte.
// If a continuation byte is invalid then error will be non-nil and the
// returned code point should be ignored. All other types of errors
// produce a non-nil error value.
func ReadCodePoint(r io.Reader) (uint32, error) {
	contLen := 0
	var offset uint32 = 0
	var cp uint32 = 0
	b := make([]byte, MaxBytes)

	for {
		n, err := r.Read(b[0:1])
		if err != nil {
			return 0, err
		}
		if n != 1 {
			return 0, fmt.Errorf("unexpected number of bytes read: %d", n)
		}

		if b[0]&(1<<7) == 0 { // 0xxxxxxx
			return uint32(b[0]), nil
		} else if (b[0]&0xE0)^0xC0 == 0 { // 110xxxxx
			contLen = 1
			offset = uint32(contLen * 6)
			cp = (0x1F & uint32(b[0])) << offset
			break
		} else if (b[0]&0xF0)^0xE0 == 0 { // 1110xxxx
			contLen = 2
			offset = uint32(contLen * 6)
			cp = (0xF & uint32(b[0])) << offset
			break
		} else if (b[0]&0xF8)^0xF0 == 0 { // 11110xxx
			contLen = 3
			offset = uint32(contLen * 6)
			cp = (0x7 & uint32(b[0])) << offset
			break
		}
		// else self-synchronize: read the next byte until we find a valid leading byte.
	}

	contBytes := b[1 : 1+contLen]
	n, err := r.Read(contBytes)
	if err != nil {
		return 0, err
	}
	if n != contLen {
		return 0, fmt.Errorf("unexpected number of bytes read: %d", n)
	}
	for _, byte := range contBytes {
		if (byte&0xC0)^0x80 != 0 { // 10xxxxxx
			return 0, fmt.Errorf("unexpected continuation byte: 0x%x\n", byte)
		}
		contValue := (byte & 0x3F)
		offset = offset - 6
		cp = cp | (uint32(contValue) << offset)
	}

	return cp, nil
}
