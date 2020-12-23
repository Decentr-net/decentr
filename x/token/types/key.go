package types

const (
	Denominator int64 = 1 * 10e6
)

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

// Key prefixes
var (
	StorePrefix = []byte{0x00} // prefix for keys that store balance
	StatsPrefix = []byte{0x01} // prefix for index keys
)
