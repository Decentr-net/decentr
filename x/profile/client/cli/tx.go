package cli

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/Decentr-net/decentr/x/profile/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	profileTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	profileTxCmd.AddCommand(flags.PostCommands(
		GetCmdSetPrivate(cdc),
		GetCmdSetPublic(cdc),
	)...)

	return profileTxCmd
}

// GetCmdSetPrivate is the CLI command for sending a SetPrivate transaction
func GetCmdSetPrivate(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-private [data]",
		Short: "set private profile data that you own",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			decoded, err := base64.StdEncoding.DecodeString(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetPrivate(decoded, cliCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdSetPublic is the CLI command for sending a SetPublic transaction
func GetCmdSetPublic(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-public [json]",
		Short: "set public profile data that you own",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			var public types.Public
			if err := json.Unmarshal([]byte(args[0]), &public); err != nil {
				return err
			}

			msg := types.NewMsgSetPublic(public, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
