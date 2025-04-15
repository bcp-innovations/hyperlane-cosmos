package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	ism "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security"
	post_dispatch "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch"
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

	cmd.AddCommand(
		ism.GetQueryCmd(),
		post_dispatch.GetQueryCmd(),
	)

	cmd.AddCommand(CmdMailboxes(),
		CmdMailbox(),
		CmdDelivered(),
		CmdRecipientIsm(),
		CmdVerifyDryRun(),
		CmdRegisteredISMs(),
		CmdRegisteredHooks(),
		CmdRegisteredApps(),
	)

	return cmd
}

func CmdMailboxes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailboxes",
		Short: "List all mailboxes",
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

			params := &types.QueryMailboxesRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Mailboxes(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "mailboxes")

	return cmd
}

func CmdMailbox() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailbox [id]",
		Short: "Get mailbox details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryMailboxRequest{
				Id: args[0],
			}

			res, err := queryClient.Mailbox(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdDelivered() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delivered [mailbox-id] [message-id]",
		Short: "Check if a message has been delivered",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryDeliveredRequest{
				Id:        args[0],
				MessageId: args[1],
			}

			res, err := queryClient.Delivered(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRecipientIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recipient-ism [recipient]",
		Short: "Get ISM ID for a recipient",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRecipientIsmRequest{
				Recipient: args[0],
			}

			res, err := queryClient.RecipientIsm(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdVerifyDryRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-dry-run [ism-id] [message] [metadata] [gas-limit]",
		Short: "Perform a dry run verification of a message",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryVerifyDryRunRequest{
				IsmId:    args[0],
				Message:  args[1],
				Metadata: args[2],
				GasLimit: args[3],
			}

			res, err := queryClient.VerifyDryRun(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRegisteredISMs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registered-isms",
		Short: "List all registered ISMs",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRegisteredISMs{}

			res, err := queryClient.RegisteredISMs(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRegisteredHooks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registered-hooks",
		Short: "List all registered hooks",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRegisteredHooks{}

			res, err := queryClient.RegisteredHooks(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdRegisteredApps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registered-apps",
		Short: "List all registered applications",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRegisteredApps{}

			res, err := queryClient.RegisteredApps(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
