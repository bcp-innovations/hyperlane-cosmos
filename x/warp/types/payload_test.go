package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
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

func TestParseWarpPayload_NoMetadata(t *testing.T) {
	// This test uses a Warp payload from the first $HYPER token transfer.
	// $HYPER implements a normal warp route, and so the payload does not
	// contain the metadata field.
	//
	// https://explorer.hyperlane.xyz/message/0x65385d162b193ecb93b2dd258f0f9b68279c4319fdc8ca6ffe3f111ec13ad880
	bz := common.FromHex("0x0000000000000000000000003fb137161365f273ebb8262a26569c117b6cbafb00000000000000000000000000000000000000000000000000005af3107a4000")

	payload, err := ParseWarpPayload(bz)
	require.NoError(t, err)

	// TODO(@john): Add validation here!
	fmt.Println(payload)
}

func TestParseWarpPayload_Metadata(t *testing.T) {
	// This test uses a Warp payload from the first $stHYPER token transfer.
	// $stHYPER implements a yielding warp route, and so the payload contains
	// the metadata field.
	//
	// https://explorer.hyperlane.xyz/message/0x1e0341782b65d086c867ec844901e4066452068eec95af95db7a6ad042d5b237
	bz := common.FromHex("0x000000000000000000000000172295445113a63f439a10e524c5779f470f31d90000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000002540be4000000000000000000000000000000000000000000000000000000000000000001")

	payload, err := ParseWarpPayload(bz)
	require.NoError(t, err)

	// TODO(@john): Add validation here!
	fmt.Println(payload)
}
