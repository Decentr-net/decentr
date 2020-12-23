package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/boltdb/bolt"
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
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	cerberusapi "github.com/Decentr-net/cerberus/pkg/api"

	"github.com/Decentr-net/decentr/app"
	"github.com/Decentr-net/decentr/x/community"
	"github.com/Decentr-net/decentr/x/pdv"
)

const (
	flagInvCheckPeriod = "inv-check-period"
	tokenDBFile        = "token.db"
	pdvDBFile          = "pdv.db"
	communityDBFile    = "community.db"
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
	rootCmd.AddCommand(flags.NewCompletionCmd(rootCmd, true))
	rootCmd.AddCommand(debug.Cmd(cdc))

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "start" {
			// cerberus address
			cmd.PersistentFlags().String(pdv.FlagCerberusAddr, "https://cerberus.testnet.decentr.xyz", "cerberus host address")
			viper.BindPFlag(pdv.FlagCerberusAddr, cmd.PersistentFlags().Lookup(pdv.FlagCerberusAddr))

			// moderator account
			cmd.PersistentFlags().String(community.FlagModeratorAddr, "", "community moderator account address")
			viper.BindPFlag(community.FlagModeratorAddr, cmd.PersistentFlags().Lookup(community.FlagModeratorAddr))

			break
		}
	}

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

	cerberusAddr := viper.GetString(pdv.FlagCerberusAddr)
	if cerberusAddr == "" {
		panic(fmt.Errorf("%s isn't set", pdv.FlagCerberusAddr))
	}

	communityModeratorAddr, err := sdk.AccAddressFromBech32(viper.GetString(community.FlagModeratorAddr))
	if err != nil {
		panic(fmt.Errorf("%s is invalid or empty", community.FlagModeratorAddr))
	}

	if _, err := url.ParseRequestURI(cerberusAddr); err != nil {
		panic(fmt.Errorf("failed to parse %s: %w", pdv.FlagCerberusAddr, err))
	}

	communityDB, err := bolt.Open(fmt.Sprintf("%s/data/%s", viper.GetString(cli.HomeFlag), communityDBFile), 0600, nil)
	if err != nil {
		panic(fmt.Errorf("failed to open communityDB: %w", err))
	}

	communityIndex, err := community.NewIndex(communityDB)
	if err != nil {
		panic(fmt.Errorf("failed to create community index: %w", err))
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	return app.NewDecentrApp(
		logger, db, traceStore, true, invCheckPeriod,
		cerberusapi.NewClient(cerberusAddr, secp256k1.PrivKeySecp256k1{}),
		communityModeratorAddr, communityIndex,
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
		aApp := app.NewDecentrApp(logger, db, traceStore, false, uint(1), nil, nil, nil)
		err := aApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	aApp := app.NewDecentrApp(logger, db, traceStore, true, uint(1), nil, nil, nil)
	return aApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
