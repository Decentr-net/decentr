package testutil

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

// NewAccAddress returns a sample account address
func NewAccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

// GetContext return a context with initialised store.
func GetContext(keys map[string]*sdk.KVStoreKey, tkeys map[string]*sdk.TransientStoreKey) (sdk.Context, error) {
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	for _, k := range keys {
		stateStore.MountStoreWithDB(k, sdk.StoreTypeIAVL, db)
	}
	for _, k := range tkeys {
		stateStore.MountStoreWithDB(k, sdk.StoreTypeTransient, nil)
	}

	if err := stateStore.LoadLatestVersion(); err != nil {
		return sdk.Context{}, err
	}

	return sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()), nil
}
