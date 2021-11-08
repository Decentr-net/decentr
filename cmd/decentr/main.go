package main

import (
	"os"

	"github.com/Decentr-net/decentr/config"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/tendermint/spm/cosmoscmd"

	"github.com/Decentr-net/decentr/app"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		config.AppName,
		config.AccountAddressPrefix,
		app.DefaultNodeHome,
		config.AppName,
		app.ModuleBasics,
		app.New,
	)
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
