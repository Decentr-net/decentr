package cli

import (
	"fmt"

	"github.com/Decentr-net/decentr/x/token/keeper"
	"github.com/Decentr-net/decentr/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group token queries under a subcommand
	tokenQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	tokenQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdBalance(queryRoute, cdc),
			GetCmdHistory(queryRoute, cdc),
			GetCmdPool(queryRoute, cdc),
		)...,
	)

	return tokenQueryCmd
}

// GetCmdBalance queries information about an account balance
func GetCmdBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "balance",
		Short: "Query account token balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/balance/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find balance - %s \n", key)
				return nil
			}

			var out keeper.Balance
			cdc.MustUnmarshalBinaryBare(bz, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdPool queries information about pool
func GetCmdPool(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Short: "Query community pool details",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", queryRoute), nil)
			if err != nil {
				fmt.Printf("failed to get pool details - %s\n", err)
				return nil
			}

			var out keeper.Pool
			cdc.MustUnmarshalBinaryBare(bz, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdHistory queries information about an account rewards history
func GetCmdHistory(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Query account rewards history",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/history/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find history - %s \n", key)
				return nil
			}

			var out []types.RewardDistribution
			cdc.MustUnmarshalBinaryBare(bz, &out)

			return cliCtx.PrintOutput(out)
		},
	}
}
