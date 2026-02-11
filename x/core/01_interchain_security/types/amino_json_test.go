package types_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	txv1beta1 "cosmossdk.io/api/cosmos/tx/v1beta1"
	"cosmossdk.io/x/tx/signing"
	"cosmossdk.io/x/tx/signing/aminojson"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"

	gogoproto "github.com/cosmos/gogoproto/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// TestAminoJSONSigningMsgSetRoutingIsmDomain reproduces the bug where Amino
// JSON signing fails for MsgSetRoutingIsmDomain because the nested Route
// message's field descriptors are not resolvable via gogoproto.HybridResolver.
//
// See: https://github.com/celestiaorg/celestia-app/issues/6541
func TestAminoJSONSigningMsgSetRoutingIsmDomain(t *testing.T) {
	// 1. Construct a MsgSetRoutingIsmDomain with a populated Route
	ismAddr := util.CreateMockHexAddress("routing-ism", 1)
	routeIsm := util.CreateMockHexAddress("route-ism", 1)

	msg := &types.MsgSetRoutingIsmDomain{
		IsmId: ismAddr,
		Route: types.Route{
			Ism:    routeIsm,
			Domain: 42,
		},
		Owner: "cosmos1testowner",
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
	// nested Route message and fails because its field descriptors are not
	// properly registered in the google.golang.org/protobuf registry.
	signBytes, err := handler.GetSignBytes(context.Background(), signerData, txData)
	require.NoError(t, err, "GetSignBytes should succeed for MsgSetRoutingIsmDomain")
	require.NotEmpty(t, signBytes)
}
