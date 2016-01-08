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

func ReadCodePoint(r io.Reader) (uint32, error) {
	b := make([]byte, MaxBytes)
	n, err := r.Read(b[0:1])
	if err != nil {
		return 0, err
	}
	if n != 1 {
		return 0, fmt.Errorf("unexpected number of bytes read: %d", n)
	}
	contLen := 0
	var offset uint32 = 0
	var cp uint32 = 0

	if b[0]&(1<<7) == 0 { // 0xxxxxxx
		return uint32(b[0]), nil
	} else if (b[0]&0xE0)^0xC0 == 0 { // 110xxxxx
		contLen = 1
		offset = uint32(contLen * 6)
		cp = (0x1F & uint32(b[0])) << offset
	} else if (b[0]&0xF0)^0xE0 == 0 { // 1110xxxx
		contLen = 2
		offset = uint32(contLen * 6)
		cp = (0xF & uint32(b[0])) << offset
	} else if (b[0]&0xF8)^0xF0 == 0 { // 11110xxx
		contLen = 3
		offset = uint32(contLen * 6)
		cp = (0x7 & uint32(b[0])) << offset
	} else {
		return 0, fmt.Errorf("unexpected leading byte: 0x%x\n", b[0])
	}

	contBytes := b[1 : 1+contLen]
	n, err = r.Read(contBytes)
	if err != nil {
		return 0, err
	}
	if n != contLen {
		return 0, fmt.Errorf("unexpected number of bytes read: %d", n)
	}
	for _, byte := range contBytes {
		// skip validation of continuation bytes
		contValue := (byte & 0x3F) // 10xxxxxx
		offset = offset - 6
		cp = cp | (uint32(contValue) << offset)
	}

	return cp, nil
}
