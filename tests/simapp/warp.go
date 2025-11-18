package simapp

import (
	warp "github.com/bcp-innovations/hyperlane-cosmos/x/warp"
	warpTypes "github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func (app *App) RegisterWarpApps() {

	warpApps := []warp.WarpApp{
		{TokenType: warpTypes.HYP_TOKEN_TYPE_COLLATERAL, App: &app.WarpKeeper},
		{TokenType: warpTypes.HYP_TOKEN_TYPE_SYNTHETIC, App: &app.WarpKeeper},
	}

	warp.RegisterWarpApp(app.HyperlaneKeeper, warpApps...)
}
