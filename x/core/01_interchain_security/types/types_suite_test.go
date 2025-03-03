package types_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestISMTypes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/hyperlane/%s Types Test Suite", types.SubModuleName))
}

func hyperlaneMessageFromHexString(hexString string) util.HyperlaneMessage {
	bytes, err := hex.DecodeString(hexString)
	Expect(err).To(BeNil())

	message, err := util.ParseHyperlaneMessage(bytes)
	Expect(err).To(BeNil())

	return message
}

func bytesFromHexString(hexString string) []byte {
	bytes, err := hex.DecodeString(hexString)
	Expect(err).To(BeNil())
	return bytes
}
