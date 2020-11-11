package cli

import (
	"bufio"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	"github.com/Decentr-net/decentr/x/community/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	communityTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	communityTxCmd.AddCommand(flags.PostCommands(
		GetCmdCreatePost(cdc),
		GetCmdDeletePost(cdc),
	)...)

	return communityTxCmd
}

// GetCmdCreatePost is the CLI command for sending a CreatePost transaction
func GetCmdCreatePost(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-post [text] [--title title] [--preview-image url] [--tag tag]",
		Short: "create blog post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			var (
				f   = cmd.Flags()
				err error

				title        string
				previewImage string
				tags         []string
			)

			if title, err = f.GetString("title"); err != nil {
				return err
			}

			if previewImage, err = f.GetString("preview-image"); err != nil {
				return err
			}

			if tags, err = f.GetStringArray("tag"); err != nil {
				return err
			}

			msg := types.NewMsgCreatePost(title, previewImage, args[0], tags, cliCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String("title", "", "post's title")
	cmd.Flags().String("preview-image", "", "post's preview image")
	cmd.Flags().StringArray("tag", nil, "post's tag")

	_ = cmd.MarkFlagRequired("title")

	return cmd
}

// GetCmdDeletePost is the CLI command for sending a DeletePost transaction
func GetCmdDeletePost(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-post [uuid]",
		Short: "delete blog post",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			id, err := uuid.FromString(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse uuid: %w", err)
			}

			msg := types.NewMsgDeletePost(id, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
