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
		GetCmdLikePost(cdc),
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
				category     int
			)

			if title, err = f.GetString("title"); err != nil {
				return err
			}

			if previewImage, err = f.GetString("preview-image"); err != nil {
				return err
			}

			if category, err = f.GetInt("category"); err != nil {
				return err
			}

			msg := types.NewMsgCreatePost(title, types.Category(category), previewImage, args[0], cliCtx.GetFromAddress())
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String("title", "", "post's title")
	cmd.Flags().Int("category", int(types.UndefinedCategory), "post's category")
	cmd.Flags().String("preview-image", "", "post's preview image")

	_ = cmd.MarkFlagRequired("title")
	_ = cmd.MarkFlagRequired("category")

	return cmd
}

// GetCmdDeletePost is the CLI command for sending a DeletePost transaction
func GetCmdDeletePost(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delete-post [postOwner] [postUUID]",
		Short: "delete blog post",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			postOwner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse postOwner: %w", err)
			}

			postUUID, err := uuid.FromString(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse postUUID: %w", err)
			}

			msg := types.NewMsgDeletePost(cliCtx.GetFromAddress(), postUUID, postOwner)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdLikePost is the CLI command for sending a LikePost transaction
func GetCmdLikePost(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "like-post [postOwner] [postUUID] [--weight weight] ",
		Short: "like blog post",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			var (
				f   = cmd.Flags()
				err error

				weightInt8 int8
			)

			postOwner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse postOwner: %w", err)
			}

			postUUID, err := uuid.FromString(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse postUUID: %w", err)
			}

			if weightInt8, err = f.GetInt8("weight"); err != nil {
				return fmt.Errorf("failed to parse weight: %w", err)
			}

			msg := types.NewMsgSetLike(postOwner, postUUID, cliCtx.GetFromAddress(), types.LikeWeight(weightInt8))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int8("weight", int8(types.LikeWeightUp),
		fmt.Sprintf("weight: like=%d dislike=%d delete=%d",
			types.LikeWeightUp, types.LikeWeightDown, types.LikeWeightZero))

	return cmd
}
