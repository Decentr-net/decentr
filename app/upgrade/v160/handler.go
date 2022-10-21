package v160

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

const Name = "v160"

func Handler(
	cfg module.Configurator,
	mm *module.Manager,
) func(ctx sdk.Context, _ upgrade.Plan, _ module.VersionMap) (module.VersionMap, error) {
	return func(ctx sdk.Context, _ upgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, cfg, fromVM)
	}
}
