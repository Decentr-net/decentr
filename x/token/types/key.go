package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// ModuleName is the name of the module
	ModuleName = "token"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

var Denominator = sdk.NewInt(1e6)

// Key prefixes
var (
	StorePrefix   = []byte{0x00} // prefix for keys that store balance
	DeltaPrefix   = []byte{0x01} // prefix for keys that store pdv delta between accruals
	RewardsPrefix = []byte{0x02} // prefix for keys for rewards history
	BanListPrefix = []byte{0x03} // prefix for keys that store bans
)

var (
	AccumulatedDelta = sdk.AccAddress{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
)
