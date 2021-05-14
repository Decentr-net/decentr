package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/Decentr-net/decentr/x/operations/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/codec"
	keyring "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	pdvTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"op"},
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	pdvTxCmd.AddCommand(flags.PostCommands(
		GetCmdSignPDV(cdc),
		GetCmdResetAccount(cdc),
	)...)

	return pdvTxCmd
}

func GetCmdSignPDV(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sign <file>",
		Short: "sign <file>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			kb, err := keyring.NewKeyring(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), inBuf)
			if err != nil {
				return fmt.Errorf("failed to get keyring: %w", err)
			}

			pdv, err := ioutil.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read pdv file: %w", err)
			}
			msg := append(pdv, []byte("/v1/pdv")...)

			signature, pk, err := kb.Sign(cliCtx.GetFromName(), keys.DefaultKeyPass, msg)
			if err != nil {
				return fmt.Errorf("failed to sign: %w", err)
			}

			return cliCtx.PrintOutput(struct {
				PublicKey string `json:"pubic_key"`
				Signature string `json:"signature"`
			}{
				PublicKey: hex.EncodeToString(pk.Bytes()[5:]), // cut amino codec prefix
				Signature: hex.EncodeToString(signature),
			})
		},
	}
}

// GetCmdResetAccount is the CLI command for sending a ResetAccount transaction
func GetCmdResetAccount(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "reset-account <account owner>",
		Short: "reset account",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			accountOwner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse account owner: %w", err)
			}

			msg := types.NewMsgResetAccount(cliCtx.GetFromAddress(), accountOwner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
