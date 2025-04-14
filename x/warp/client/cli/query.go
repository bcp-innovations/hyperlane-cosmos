package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group query queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdTokens())
	cmd.AddCommand(CmdToken())
	cmd.AddCommand(CmdBridgedSupply())
	cmd.AddCommand(CmdRemoteRouters())
	cmd.AddCommand(CmdQuoteRemoteTransfer())

	return cmd
}

func CmdTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "List all tokens",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTokensRequest{}

			res, err := queryClient.Tokens(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [id]",
		Short: "Get token details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTokenRequest{
				Id: args[0],
			}

			res, err := queryClient.Token(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdBridgedSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridged-supply [id]",
		Short: "Query bridged supply of a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryBridgedSupplyRequest{
				Id: args[0],
			}

			res, err := queryClient.BridgedSupply(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRemoteRouters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote-routers [id]",
		Short: "Query remote routers for a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryRemoteRoutersRequest{
				Id:         args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.RemoteRouters(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "remote-routers")

	return cmd
}

func CmdQuoteRemoteTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote-transfer [id] [destination-domain]",
		Short: "Quote gas payment for a remote transfer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryQuoteRemoteTransferRequest{
				Id:                args[0],
				DestinationDomain: args[1],
			}

			res, err := queryClient.QuoteRemoteTransfer(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
