package cli

import (
	"fmt"

	"github.com/Decentr-net/decentr/x/operations/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group profile queries under a subcommand
	profileQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	profileQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdMinGasPrice(cdc),
		)...,
	)

	return profileQueryCmd
}

func GetCmdMinGasPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "min-gas-price",
		Short: "Query the current min gas price value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMinGasPrice)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var mgp sdk.DecCoin
			if err := cdc.UnmarshalJSON(res, &mgp); err != nil {
				return err
			}

			return cliCtx.PrintOutput(mgp)
		},
	}
}
