package utf8

import (
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []rune{
		'本', // 0x672c
	}
	for _, testValue := range tests {
		r := strings.NewReader(string(testValue))
		cp, err := ReadCodePoint(r)
		if err != nil {
			t.Error("unexpected error: ", err)
		}
		if cp != uint32(testValue) {
			t.Errorf("expected 0x%x, got 0x%x", testValue, cp)
		}
	}
}
