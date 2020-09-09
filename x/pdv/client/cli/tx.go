package cli

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	keyring "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	pdvTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	pdvTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreatePDV(cdc),
		GetCmdSignPDV(cdc),
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
			digest := sha256.Sum256(append(pdv, []byte("/v1/pdv")...))

			signature, pk, err := kb.Sign(cliCtx.GetFromName(), keys.DefaultKeyPass, digest[:])
			if err != nil {
				return fmt.Errorf("failed to sign: %w", err)
			}

			return cliCtx.PrintOutput(struct {
				PublicKey string `json:"pubic_key"`
				Signature string `json:"signature"`
			}{
				PublicKey: hex.EncodeToString(pk.Bytes()),
				Signature: hex.EncodeToString(signature),
			})
		},
	}
}

// GetCmdCreatePDV is the CLI command for sending a CreatePDV transaction
func GetCmdCreatePDV(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create [pdv]",
		Short: "create PDV",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			if hex.EncodeToString(cliCtx.GetFromAddress()) != strings.Split(args[0], "-")[0] { // Checks if the the msg sender is the same as the current owner
				return fmt.Errorf("invalid owner")
			}

			caddr, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/cerberus-addr", types.QuerierRoute), nil)
			if err != nil {
				return fmt.Errorf("failed to get cerberus addr: %w", err)
			}

			exists, err := cerberusapi.NewClient(string(caddr), secp256k1.PrivKeySecp256k1{}).DoesPDVExist(cmd.Context(), args[0])
			if err != nil {
				return fmt.Errorf("failed to check pdv existence: %w", err)
			}

			if !exists {
				return fmt.Errorf("pdv does not exist")
			}

			msg := types.NewMsgCreatePDV(time.Now().UTC(), args[0], types.PDVTypeCookie, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
