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
		CmdAggregationHooks(),
		CmdAggregationHook(),
		CmdPausableHooks(),
		CmdPausableHook(),
		CmdRateLimitedHooks(),
		CmdRateLimitedHook(),
		CmdTokenRateLimits(),
		CmdTokenRateLimit(),
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

func CmdAggregationHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregation-hooks",
		Short: "List all aggregation hooks",
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

			params := &types.QueryAggregationHooksRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AggregationHooks(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "aggregation-hooks")

	return cmd
}

func CmdAggregationHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregation-hook [id]",
		Short: "Get details for a specific aggregation hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAggregationHookRequest{
				Id: args[0],
			}

			res, err := queryClient.AggregationHook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdPausableHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pausable-hooks",
		Short: "List all pausable hooks",
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

			params := &types.QueryPausableHooksRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.PausableHooks(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "pausable-hooks")

	return cmd
}

func CmdPausableHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pausable-hook [id]",
		Short: "Get details for a specific pausable hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryPausableHookRequest{
				Id: args[0],
			}

			res, err := queryClient.PausableHook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRateLimitedHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate-limited-hooks",
		Short: "List all rate limited hooks",
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

			params := &types.QueryRateLimitedHooksRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RateLimitedHooks(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "rate-limited-hooks")

	return cmd
}

func CmdRateLimitedHook() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rate-limited-hook [id]",
		Short: "Get details for a specific rate limited hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRateLimitedHookRequest{
				Id: args[0],
			}

			res, err := queryClient.RateLimitedHook(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdTokenRateLimits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-rate-limits [hook-id]",
		Short: "List token rate limits for a rate limited hook",
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

			params := &types.QueryTokenRateLimitsRequest{
				HookId:     args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.TokenRateLimits(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "token-rate-limits")

	return cmd
}

func CmdTokenRateLimit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token-rate-limit [hook-id] [token-id]",
		Short: "Get a token rate limit for a rate limited hook",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryTokenRateLimitRequest{
				HookId:  args[0],
				TokenId: args[1],
			}

			res, err := queryClient.TokenRateLimit(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
