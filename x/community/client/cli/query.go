package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/Decentr-net/decentr/x/community/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group community queries under a subcommand
	communityQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	communityQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdPost(queryRoute, cdc),
			GetCmdModerators(queryRoute, cdc),
			GetCmdUsersPosts(queryRoute, cdc),
			GetCmdFollowee(queryRoute, cdc),
		)...,
	)

	return communityQueryCmd
}

// GetCmdPost queries exact post
func GetCmdPost(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "post <owner> <uuid>",
		Short: "Query post",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/post/%s/%s", queryRoute, args[0], args[1]), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdFollowee queries users followee
func GetCmdFollowee(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "followee <owner>",
		Short: "Query user's followee",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/followees/%s", queryRoute, args[0]), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdUsersPosts queries users posts
func GetCmdUsersPosts(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user-posts <owner> [--from uuid] [--limit int]",
		Short: "Query user's posts",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				f     = cmd.Flags()
				from  string
				limit int

				err error
			)

			if from, err = f.GetString("from-uuid"); err != nil {
				return err
			}

			if limit, err = f.GetInt("limit"); err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/user/%s/%s/%d", queryRoute, args[0], from, limit), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String("from-uuid", "", "list from uuid")
	cmd.Flags().Int("limit", 20, "limit")

	return cmd
}

// GetCmdModerators queries for the community moderators
func GetCmdModerators(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "moderators",
		Short: "Returns moderators addresses",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/moderators", queryRoute), nil)
			if err != nil {
				fmt.Printf("failed to get moderators - %s \n", err.Error())
				return nil
			}
			return cliCtx.PrintOutput(string(res))
		},
	}
}
