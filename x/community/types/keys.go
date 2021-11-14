package types

const (
	// ModuleName defines the module name
	ModuleName = "community"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	PostPrefix      = []byte{0x01} // prefix for keys that store posts
	LikePrefix      = []byte{0x02} // prefix for keys that store likes
	FollowingPrefix = []byte{0x03} // prefix for keys that store followings
)
