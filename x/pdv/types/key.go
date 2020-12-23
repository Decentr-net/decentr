package types

const (
	// ModuleName is the name of the module
	ModuleName = "pdv"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName

	FlagCerberusAddr = "cerberus-addr"
)

// Key prefixes
var (
	StorePrefix = []byte{0x00} // prefix for keys that store balance
	IndexPrefix = []byte{0x01} // prefix for index keys
)
