package util_test

import (
	"encoding/hex"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWarpKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "./util Test Suite")
}

func MustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	Expect(err).To(BeNil())
	return b
}
