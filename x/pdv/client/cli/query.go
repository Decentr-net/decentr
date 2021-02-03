package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/Decentr-net/decentr/x/pdv/types"
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
			GetCmdCerberusAddr(queryRoute, cdc),
		)...,
	)

	return pdvQueryCmd
}

// GetCmdCerberusAddr queries for the cerberus-addr flag
func GetCmdCerberusAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "cerberus-addr",
		Short: "Returns current cerberus address",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/cerberus-addr", queryRoute), nil)
			if err != nil {
				fmt.Printf("failed to get cerberus addr - %s \n", err.Error())
				return nil
			}
			return cliCtx.PrintOutput(string(res))
		},
	}
}
