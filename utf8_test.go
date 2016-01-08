package utf8

import (
	"io"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []string{
		"本",
		"Hello, 世界",
	}
	for _, test := range tests {
		s := ""
		r := strings.NewReader(test)
		for {
			cp, err := ReadCodePoint(r)
			if err == io.EOF {
				break
			} else if err != nil {
				t.Error("unexpected error: ", err)
			}
			s = s + string(rune(cp))
		}
		if s != test {
			t.Error("")
		}
	}
}
