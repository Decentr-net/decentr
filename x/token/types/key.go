package types

const (
	Denominator int64 = 1 * 10e6
	MaxSupply   int64 = 1 * 10e9 * Denominator
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
