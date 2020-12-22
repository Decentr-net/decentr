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
			GetCmdOwner(queryRoute, cdc),
			GetCmdShow(queryRoute, cdc),
			GetCmdList(queryRoute, cdc),
			GetCmdCerberusAddr(queryRoute, cdc),
		)...,
	)

	return pdvQueryCmd
}

// GetCmdPDV queries PDV owner
func GetCmdOwner(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "owner <address>",
		Short: "Query PDV owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/owner/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find PDV - %s owner \n", key)
				return nil
			}
			return cliCtx.PrintOutput(string(res))
		},
	}
}

// GetCmdShow queries PDV full data unencrypted data
func GetCmdShow(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "show <address>",
		Short: "Query PDV meta data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/show/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find PDV - %s \n", key)
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdShow queries PDV full data unencrypted data
func GetCmdList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list <owner> [from] [limit]",
		Short: "Query list of PDVs meta data. Default from is 0001-01-01T00:00:00Z and limit is 20",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			owner, limit, from := args[0], "", ""
			if len(args) > 1 {
				from = args[1]
				if len(args) > 2 {
					limit = args[2]
				}
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/list/%s/%s/%s", queryRoute, owner, from, limit), nil)
			if err != nil {
				fmt.Printf("could not list PDV - %s \n", err.Error())
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}
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
