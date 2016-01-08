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
		runeValue, err := ReadFrom(r)
		if err != nil {
			t.Error("unexpected error: ", err)
		}
		if runeValue != testValue {
			t.Errorf("expected 0x%x, got 0x%x", testValue, runeValue)
		}
	}
}
