package types

import (
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
			if isZeroPadded(p.bz) != p.ok {
				t.Fail()
			}
		})
	}
}
