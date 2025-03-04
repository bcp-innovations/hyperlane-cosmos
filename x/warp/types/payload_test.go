package types

import (
	"bytes"
	"testing"
)

func TestIsZeroPadded(t *testing.T) {
	type pair struct {
		bz []byte
		ok bool
	}
	for _, p := range []pair{
		{[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, true},
		{[]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, false},
		{[]byte{1}, false},
		{[]byte{}, false},
	} {
		t.Run(string(p.bz), func(t *testing.T) {
			t.Log(bytes.HasPrefix(p.bz, make([]byte, 0, 12)))
			if isZeroPadded(p.bz) != p.ok {
				t.Fail()
			}
		})
	}
}
