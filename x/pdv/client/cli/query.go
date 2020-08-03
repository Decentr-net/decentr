package cli

import (
	"fmt"

	"github.com/Decentr-net/decentr/x/pdv/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group pdv queries under a subcommand
	pdvQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	pdvQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdOwner(queryRoute, cdc),
			GetCmdShow(queryRoute, cdc),
		)...,
	)

	return pdvQueryCmd
}

// GetCmdPDV queries PDV owner
func GetCmdOwner(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "owner",
		Short: "Query PDV owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/owner/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find PDV - %s \n", key)
				return nil
			}
			return cliCtx.PrintOutput(string(res))
		},
	}
}

// GetCmdShow queries PDV full data unencrypted data
func GetCmdShow(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Query PDV unencrypted data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/show/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find PDV - %s \n", key)
				return nil
			}
			return cliCtx.PrintOutput(string(res))
		},
	}
}
