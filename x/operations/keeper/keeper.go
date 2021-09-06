package keeper

import (
	"github.com/Decentr-net/decentr/x/operations/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc      *codec.Codec // The wire codec for binary encoding/decoding.

	paramSubspace params.Subspace
}

// NewKeeper creates new instances of the PDV Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, paramSpace params.Subspace) Keeper {
	ps := paramSpace.WithKeyTable(types.ParamKeyTable())

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSubspace: ps,
	}
}
