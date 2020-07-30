package cli

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/profile/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
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
			GetCmdPrivate(queryRoute, cdc),
			GetCmdPublic(queryRoute, cdc),
		)...,
	)

	return profileQueryCmd
}

// GetCmdPrivate queries information about a private profile
func GetCmdPrivate(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "private",
		Short: "Query private profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/private/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find private profile - %s \n", key)
				return nil
			}
			return cliCtx.PrintOutput(string(res))
		},
	}
}

// GetCmdPublic queries information about a public profile
func GetCmdPublic(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "public",
		Short: "Query public profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/public/%s", queryRoute, key), nil)
			if err != nil {
				fmt.Printf("could not find public profile - %s \n", key)
				return nil
			}

			var out types.Public
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
