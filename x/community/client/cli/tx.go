package cli

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/types"
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewCreatePostCmd(),
		NewDeletePostCmd(),
		NewSetLikeCmd(),
		NewFollowCmd(),
		NewUnfollowCmd(),
	)

	return txCmd
}

func NewCreatePostCmd() *cobra.Command {
	categories := make(map[string]int, len(types.Category_value))
	categoriesList := make([]string, 0, len(types.Category_value))
	for k, v := range types.Category_value {
		category := strings.TrimLeft(k, "CATEGORY_")
		categories[category] = int(v)
		categoriesList = append(categoriesList, category)
	}

	cmd := &cobra.Command{
		Use:   "create-post [text] --title [title] --category [category] <--preview-image [url]>",
		Short: "Create new post",
		Long:  fmt.Sprintf("[category] can be on of %s", strings.Join(categoriesList, ", ")),
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			f := cmd.Flags()

			title, err := f.GetString("title")
			if err != nil {
				return fmt.Errorf("invalid title: %w", err)
			}

			var (
				category int
				ok       bool
			)
			if v, err := f.GetString("category"); err != nil {
				return fmt.Errorf("invalid category: %w", err)
			} else if category, ok = categories[v]; !ok {
				return fmt.Errorf("invalid category: must be one of %s", strings.Join(categoriesList, ", "))
			}

			previewImage, err := f.GetString("preview-image")
			if err != nil {
				return fmt.Errorf("invalid preview-image: %w", err)
			}

			msg := types.NewMsgCreatePost(
				title,
				types.Category(category),
				previewImage,
				args[0],
				clientCtx.GetFromAddress(),
			)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().String("title", "", "post's title")
	cmd.Flags().String("category", "", "post's category")
	cmd.Flags().String("preview-image", "", "post's preview image")

	_ = cmd.MarkFlagRequired("title")    // nolint
	_ = cmd.MarkFlagRequired("category") // nolint

	return cmd
}

func NewDeletePostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-post [owner] [uuid]",
		Short: "Delete a post by [owner] and [uuid]",
		Long:  "[owner] can be omit and used from --from flag",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			owner := clientCtx.GetFromAddress()

			if len(args) > 1 {
				owner, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return fmt.Errorf("invalid owner: %w", err)
				}
			}

			id, err := uuid.FromString(args[len(args)-1])
			if err != nil {
				return fmt.Errorf("invalid uuid: %w", err)
			}

			msg := types.NewMsgDeletePost(
				clientCtx.GetFromAddress(),
				owner,
				id,
			)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewSetLikeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "like-post [owner] [uuid] <--weight [weight]>",
		Short: "Like a post by [owner] and [uuid] with optionally set [weight]",
		Args:  cobra.RangeArgs(2, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			id, err := uuid.FromString(args[1])
			if err != nil {
				return fmt.Errorf("invalid uuid: %w", err)
			}

			weight, err := cmd.Flags().GetInt("weight")
			if err != nil {
				return fmt.Errorf("invalid weight: %w", err)
			}

			msg := types.NewMsgSetLike(
				owner,
				id,
				clientCtx.GetFromAddress(),
				types.LikeWeight(weight),
			)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Int(
		"weight",
		int(types.LikeWeight_LIKE_WEIGHT_UP),
		"Weight for like operation. Use 1 to like post, -1 to dislike and 0 to remove your reaction",
	)

	return cmd
}

func NewFollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "follow [address]",
		Short: "Follow an account by [address]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			msg := types.NewMsgFollow(
				clientCtx.GetFromAddress(),
				address,
			)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func NewUnfollowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unfollow [address]",
		Short: "Unfollow an account by [address]",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().Empty() {
				return fmt.Errorf("--from flag should be specified")
			}

			address, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			msg := types.NewMsgUnfollow(
				clientCtx.GetFromAddress(),
				address,
			)
			if err := msg.ValidateBasic(); err != nil {
				return fmt.Errorf("invalid msg: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
