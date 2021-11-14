package cli

import (
	"fmt"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Decentr-net/decentr/x/community/types"
)

func GetQueryCmd(_ string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewModeratorsCmd(),
		NewGetPostCmd(),
		NewListUserPostsCmd(),
		NewListFollowedCmd(),
	)

	return cmd
}

func NewModeratorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moderators",
		Short: "Query the moderators list",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.Moderators(cmd.Context(), &types.ModeratorsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewGetPostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post [owner] [uuid]",
		Short: "Query post by [owner] and [uuid]",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			id, err := uuid.FromString(args[1])
			if err != nil {
				return fmt.Errorf("invalid uuid: %w", err)
			}

			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.GetPost(cmd.Context(), &types.GetPostRequest{
				PostOwner: owner,
				PostUuid:  id.String(),
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewListUserPostsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-posts [owner] <offset> <limit>",
		Short: "Lists post by [owner] with optional <limit>/<offset> pagination",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			var offset, limit uint64
			if len(args) > 1 {
				if offset, err = strconv.ParseUint(args[1], 10, 64); err != nil {
					return fmt.Errorf("invalid offset: %w", err)
				}
			}

			if len(args) > 2 {
				if limit, err = strconv.ParseUint(args[2], 10, 64); err != nil {
					return fmt.Errorf("invalid offset: %w", err)
				}
			}

			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.ListUserPosts(cmd.Context(), &types.ListUserPostsRequest{
				Owner: owner,
				Pagination: query.PageRequest{
					Offset: offset,
					Limit:  limit,
				},
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func NewListFollowedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "followed [owner] <offset> <limit>",
		Short: "Lists followed by [owner] account addresses with optional <limit>/<offset> pagination",
		Args:  cobra.RangeArgs(1, 3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			owner, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid owner: %w", err)
			}

			var offset, limit uint64
			if len(args) > 1 {
				if offset, err = strconv.ParseUint(args[1], 10, 64); err != nil {
					return fmt.Errorf("invalid offset: %w", err)
				}
			}

			if len(args) > 2 {
				if limit, err = strconv.ParseUint(args[2], 10, 64); err != nil {
					return fmt.Errorf("invalid offset: %w", err)
				}
			}

			queryClient := types.NewQueryClient(clientCtx)

			out, err := queryClient.ListFollowed(cmd.Context(), &types.ListFollowedRequest{
				Owner: owner,
				Pagination: query.PageRequest{
					Offset: offset,
					Limit:  limit,
				},
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(out)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
