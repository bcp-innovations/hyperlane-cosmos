package types

import "github.com/bcp-innovations/hyperlane-cosmos/util"

type RoutingISM interface {
	Route(message util.HyperlaneMessage) (util.HexAddress, error)
}
