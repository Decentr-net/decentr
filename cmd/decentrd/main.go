package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/Decentr-net/decentr/x/community"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/Decentr-net/decentr/app"
)

const (
	flagInvCheckPeriod = "inv-check-period"
)

var invCheckPeriod uint

func main() {
	cdc := app.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:               "decentrd",
		Short:             "Decentr Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	rootCmd.AddCommand(genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome))
	rootCmd.AddCommand(genutilcli.MigrateGenesisCmd(ctx, cdc))
	rootCmd.AddCommand(
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{},
			auth.GenesisAccountIterator{}, app.DefaultNodeHome, app.DefaultCLIHome,
		),
	)
	rootCmd.AddCommand(genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics))
	rootCmd.AddCommand(AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(AddGenesisCommunityModeratorsCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome))
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(debug.Cmd(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "AU", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	var cache sdk.MultiStorePersistentCache

	if viper.GetBool(server.FlagInterBlockCache) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	return app.NewDecentrApp(
		logger, db, traceStore, true, invCheckPeriod,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
		baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
		baseapp.SetInterBlockCache(cache),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		aApp := app.NewDecentrApp(logger, db, traceStore, false, uint(1))
		err := aApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	aApp := app.NewDecentrApp(logger, db, traceStore, true, uint(1))
	return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}

// AddGenesisCommunityModeratorsCmd returns add-genesis-community-moderators cobra Command.
func AddGenesisCommunityModeratorsCmd(
	ctx *server.Context, cdc *codec.Codec, defaultNodeHome, defaultClientHome string,
) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add-genesis-community-moderators [moderator][,[moderator]]",
		Short: "Add a genesis community moderators to genesis.json",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			moderators := strings.Split(args[0], ",")
			for _, moderator := range moderators {
				if _, err := sdk.AccAddressFromBech32(moderator); err != nil {
					return fmt.Errorf("failed to parse moderators: %w", err)
				}
			}

			genFile := config.GenesisFile()
			appState, genDoc, err := genutil.GenesisStateFromGenFile(cdc, genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			communityGenState := community.GetGenesisStateFromAppState(cdc, appState)
			communityGenState.Moderators = moderators

			communityGenStateBz, err := cdc.MarshalJSON(communityGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal community genesis state: %w", err)
			}

			appState[community.ModuleName] = communityGenStateBz

			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")

	return cmd
}
