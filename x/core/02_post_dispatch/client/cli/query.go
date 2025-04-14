package cli

import (
	"fmt"
	"strings"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {

	// Group query queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "hooks",
		Short:                      fmt.Sprintf("Querying commands for the %s module", strings.Replace(types.SubModuleName, "_", "-", 1)),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdIgps(),
		CmdIgp(),
		CmdDestinationGasConfigs(),
		CmdQuoteGasPayment(),
		CmdMerkleTreeHooks(),
		CmdMerkleTreeHook(),
		CmdNoopHooks(),
		CmdNoopHook(),
	)
	return cmd
}

func CmdIgps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "igps",
		Short: "List all interchain gas paymasters",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryIgpsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Igps(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "igps")

	return cmd
}

func CmdIgp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "igp [id]",
		Short: "Get details for a specific interchain gas paymaster",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryIgpRequest{
				Id: args[0],
			}

			res, err := queryClient.Igp(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdDestinationGasConfigs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destination-gas-configs [id]",
		Short: "List destination gas configs for an IGP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryDestinationGasConfigsRequest{
				Id:         args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.DestinationGasConfigs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "destination-gas-configs")

	return cmd
}

func CmdQuoteGasPayment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote-gas-payment [igp-id] [destination-domain] [gas-limit]",
		Short: "Quote gas payment for a transaction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryQuoteGasPaymentRequest{
				IgpId:             args[0],
				DestinationDomain: args[1],
				GasLimit:          args[2],
			}

			res, err := queryClient.QuoteGasPayment(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdMerkleTreeHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merkle-tree-hooks",
		Short: "List all merkle tree hooks",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryMerkleTreeHooksRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.MerkleTreeHooks(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "merkle-tree-hooks")

	return cmd
}

func CmdMerkleTreeHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merkle-tree-hook [id]",
		Short: "Get details for a specific merkle tree hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryMerkleTreeHookRequest{
				Id: args[0],
			}

			res, err := queryClient.MerkleTreeHook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdNoopHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "noop-hooks",
		Short: "List all noop hooks",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryNoopHooksRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.NoopHooks(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "noop-hooks")

	return cmd
}

func CmdNoopHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "noop-hook [id]",
		Short: "Get details for a specific noop hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryNoopHookRequest{
				Id: args[0],
			}

			res, err := queryClient.NoopHook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
