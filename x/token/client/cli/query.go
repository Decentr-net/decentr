package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
)

func GetQueryCmd(_ string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewBalanceCmd(),
		NewPoolCmd(),
	)

	return cmd
}

func NewBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [address]",
		Short: "Query account token balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address: %w", err)
			}
			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.Balance(cmd.Context(), &types.BalanceRequest{
				Address: address,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewPoolCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Query pool balance",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.Pool(cmd.Context(), &types.PoolRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
