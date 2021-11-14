package types

const (
	// ModuleName defines the module name
	ModuleName = "token"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	CalcPrefix = []byte{0x00}

	BalancePrefix      = []byte{0x01} // prefix for keys that store balance
	BalanceDeltaPrefix = []byte{0x02} // prefix for keys that store pdv delta between accruals
	BanListPrefix      = []byte{0x03} // prefix for keys that store bans

	AccumulatedDeltaKey = []byte{0x01}
)
