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
	PostPrefix           = []byte{0x00} // prefix for keys that store posts
	LikePrefix           = []byte{0x01} // prefix for keys that store likes
	IndexCreatedAtPrefix = []byte{0x02} // prefix for created_at index keys
	IndexPopularPrefix   = []byte{0x03} // prefix for popular index keys
	IndexUserLikesPrefix = []byte{0x04} // prefix for user likes index keys
	FollowersPrefix      = []byte{0x05} // prefix for store followers
)
