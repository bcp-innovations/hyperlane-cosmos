package types_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	txv1beta1 "cosmossdk.io/api/cosmos/tx/v1beta1"
	"cosmossdk.io/math"
	"cosmossdk.io/x/tx/signing"
	"cosmossdk.io/x/tx/signing/aminojson"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// TestAminoJSONSigningMsgEnrollRemoteRouter reproduces the bug where Amino
// JSON signing fails for MsgEnrollRemoteRouter because the nested RemoteRouter
// message's field descriptors are not resolvable via gogoproto.HybridResolver.
//
// See: https://github.com/bcp-innovations/hyperlane-cosmos/issues/156
func TestAminoJSONSigningMsgEnrollRemoteRouter(t *testing.T) {
	// 1. Construct a MsgEnrollRemoteRouter with a populated RemoteRouter
	tokenId := util.CreateMockHexAddress("token", 1)

	msg := &types.MsgEnrollRemoteRouter{
		Owner:   "cosmos1testowner",
		TokenId: tokenId,
		RemoteRouter: &types.RemoteRouter{
			ReceiverDomain:   42,
			ReceiverContract: "0x0000000000000000000000000000000000000001",
			Gas:              math.NewInt(100000),
		},
	}

	// 2. Marshal with gogo proto
	msgBytes, err := gogoproto.Marshal(msg)
	require.NoError(t, err)

	// 3. Wrap as an anypb.Any
	anyMsg := &anypb.Any{
		TypeUrl: "/" + gogoproto.MessageName(msg),
		Value:   msgBytes,
	}

	// 4. Build a TxBody and marshal to bodyBytes
	body := &txv1beta1.TxBody{
		Messages: []*anypb.Any{anyMsg},
	}
	bodyBytes, err := proto.Marshal(body)
	require.NoError(t, err)

	// 5. Build a minimal AuthInfo with a Fee and marshal to authInfoBytes
	authInfo := &txv1beta1.AuthInfo{
		Fee: &txv1beta1.Fee{
			GasLimit: 200000,
		},
	}
	authInfoBytes, err := proto.Marshal(authInfo)
	require.NoError(t, err)

	// 6. Create an amino JSON sign mode handler with default options
	// (uses gogoproto.HybridResolver, same as the real app)
	handler := aminojson.NewSignModeHandler(aminojson.SignModeHandlerOptions{})

	signerData := signing.SignerData{
		Address:       "cosmos1testowner",
		ChainID:       "test-chain",
		AccountNumber: 1,
		Sequence:      0,
	}

	txData := signing.TxData{
		Body:          body,
		AuthInfo:      authInfo,
		BodyBytes:     bodyBytes,
		AuthInfoBytes: authInfoBytes,
	}

	// 7. Call GetSignBytes — this is where the bug manifests.
	// The handler internally calls RejectUnknownFields which descends into the
	// nested RemoteRouter message and fails because its field descriptors are not
	// properly registered in the google.golang.org/protobuf registry.
	signBytes, err := handler.GetSignBytes(context.Background(), signerData, txData)
	require.NoError(t, err, "GetSignBytes should succeed for MsgEnrollRemoteRouter")
	require.NotEmpty(t, signBytes)
}
