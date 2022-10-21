package keeper

import (
	"fmt"

	"github.com/Decentr-net/decentr/config"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Decentr-net/decentr/x/operations/types"
)

type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey sdk.StoreKey

	paramSpace paramtypes.Subspace
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
) *Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		paramSpace: paramSpace,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

func (k Keeper) GetMinGasPrice(ctx sdk.Context) sdk.DecCoin {
	params := types.Params{
		MinGasPrice: sdk.NewDecCoin(config.DefaultBondDenom, sdk.ZeroInt()),
	}
	k.paramSpace.GetParamSetIfExists(ctx, &params)
	return params.MinGasPrice
}

func (k Keeper) IsSupervisor(ctx sdk.Context, address sdk.AccAddress) bool {
	for _, v := range k.GetParams(ctx).Supervisors {
		supervisor, _ := sdk.AccAddressFromBech32(v)
		if address.Equals(supervisor) {
			return true
		}
	}

	return false
}
