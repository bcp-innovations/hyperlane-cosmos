package cli

import (
	"fmt"

	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group query queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.SubModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.SubModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdIsms(),
		CmdIsm(),
		CmdAnnouncedStorageLocations(),
		CmdLatestAnnouncedStorageLocation(),
	)

	return cmd
}

func CmdIsms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isms",
		Short: "List all ISMs",
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

			params := &types.QueryIsmsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Isms(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "isms")

	return cmd
}

func CmdIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ism [id]",
		Short: "Get ISM details by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryIsmRequest{
				Id: args[0],
			}

			res, err := queryClient.Ism(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdAnnouncedStorageLocations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announced-storage-locations [mailbox-id] [validator-address]",
		Short: "Query announced storage locations for a validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAnnouncedStorageLocationsRequest{
				MailboxId:        args[0],
				ValidatorAddress: args[1],
			}

			res, err := queryClient.AnnouncedStorageLocations(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdLatestAnnouncedStorageLocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-announced-storage-location [mailbox-id] [validator-address]",
		Short: "Query the latest announced storage location for a validator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryLatestAnnouncedStorageLocationRequest{
				MailboxId:        args[0],
				ValidatorAddress: args[1],
			}

			res, err := queryClient.LatestAnnouncedStorageLocation(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
