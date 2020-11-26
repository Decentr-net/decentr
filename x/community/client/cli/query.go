package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/Decentr-net/decentr/x/community/keeper"
	"github.com/Decentr-net/decentr/x/community/types"
)

const (
	dayInterval   = "day"
	weekInterval  = "week"
	monthInterval = "month"
)

var intervals = map[string]keeper.Interval{
	dayInterval:   keeper.DayInterval,
	weekInterval:  keeper.WeekInterval,
	monthInterval: keeper.MonthInterval,
}

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
			GetCmdUsersPosts(queryRoute, cdc),
			GetCmdPopularPostsList(queryRoute, cdc),
			GetCmdPostsList(queryRoute, cdc),
		)...,
	)

	return communityQueryCmd
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

			if from, err = f.GetString("from"); err != nil {
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

	cmd.Flags().String("from", "", "list from uuid")
	cmd.Flags().Int("limit", 20, "limit")

	return cmd
}

// GetCmdPopularPostsList queries popular posts
func GetCmdPopularPostsList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "popular-posts [--from-owner owner --from-uuid uuid] [--category string] [--limit int] [--interval day/week/month]",
		Short: "Query popular posts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				f         = cmd.Flags()
				fromOwner string
				fromUUID  string
				category  int
				limit     int
				interval  string

				err error
			)

			if fromOwner, err = f.GetString("from-owner"); err != nil {
				return err
			}
			if fromUUID, err = f.GetString("from-uuid"); err != nil {
				return err
			}
			if category, err = f.GetInt("category"); err != nil {
				return err
			}
			if limit, err = f.GetInt("limit"); err != nil {
				return err
			}
			if interval, err = f.GetString("interval"); err != nil {
				return err
			}

			qPath := fmt.Sprintf("custom/%s/popular/%s/%s/%d/%d/%d", queryRoute, fromOwner, fromUUID, limit, category, intervals[interval])
			res, _, err := cliCtx.QueryWithData(qPath, nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String("from-owner", "", "list from post (owner part)")
	cmd.Flags().String("from-uuid", "", "list from post (uuid part)")
	cmd.Flags().String("interval", "month", "interval for searching (day, week, month)")
	cmd.Flags().Int("category", int(types.UndefinedCategory), "post's category")
	cmd.Flags().Int("limit", 20, "limit")

	return cmd
}

// GetCmdPostsList queries the latest posts
func GetCmdPostsList(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "posts [--from-owner owner --from-uuid uuid] [--category string] [--limit int]",
		Short: "Query the latest posts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				f         = cmd.Flags()
				fromOwner string
				fromUUID  string
				category  int
				limit     int

				err error
			)

			if fromOwner, err = f.GetString("from-owner"); err != nil {
				return err
			}
			if fromUUID, err = f.GetString("from-uuid"); err != nil {
				return err
			}

			if category, err = f.GetInt("category"); err != nil {
				return err
			}

			if limit, err = f.GetInt("limit"); err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/posts/%s/%s/%d/%d", queryRoute, fromOwner, fromUUID, limit, category), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String("from-owner", "", "list from post (owner part)")
	cmd.Flags().String("from-uuid", "", "list from post (uuid part)")
	cmd.Flags().Int("category", int(types.UndefinedCategory), "post's category")
	cmd.Flags().Int("limit", 20, "limit")

	return cmd
}
