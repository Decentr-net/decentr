package cli

import (
	goctx "context"
	"fmt"
	"strconv"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto/secp256k1"

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
			GetCmdMeta(queryRoute, cdc),
			GetCmdCerberusAddr(queryRoute, cdc),
		)...,
	)

	return pdvQueryCmd
}

// GetCmdMeta queries for the pdv's details.
func GetCmdMeta(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "meta <id> --from [account]",
		Short: "Returns pdv details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Printf("invalid id: %s\n", err.Error())
				return nil
			}

			cerberusAddr, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/cerberus-addr", queryRoute), nil)
			if err != nil {
				fmt.Printf("failed to get cerberus addr - %s \n", err.Error())
				return nil
			}

			cerberus := cerberusapi.NewClient(string(cerberusAddr), secp256k1.PrivKeySecp256k1{})
			meta, err := cerberus.GetPDVMeta(goctx.Background(), cliCtx.FromAddress.String(), id)
			if err != nil {
				fmt.Printf("failed to get pdv details - %s \n", err.Error())
				return nil
			}

			return cliCtx.PrintOutput(meta)
		},
	}
}

// GetCmdCerberusAddr queries for the cerberus-addr flag.
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
