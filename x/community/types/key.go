package types

const (
	// ModuleName is the name of the module
	ModuleName = "community"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

// Key prefixes
var (
	PostPrefix = []byte{0x00} // prefix for keys that store posts
	LikePrefix = []byte{0x01} // prefix for keys that store likes
)
