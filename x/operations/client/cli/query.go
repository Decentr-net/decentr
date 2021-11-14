package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/operations/types"
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
		NewMinGasPriceCmd(),
		NewIsAccountBannedCmd(),
	)

	return cmd
}

func NewMinGasPriceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "min-gas-price",
		Short: "Query the current min gas price value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.MinGasPrice(cmd.Context(), &types.MinGasPriceRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewIsAccountBannedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-account-banned [address]",
		Short: "Query if the account banned",
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

			out, err := queryClient.IsAccountBanned(cmd.Context(), &types.IsAccountBannedRequest{
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
