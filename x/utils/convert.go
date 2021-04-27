package utils

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/token/types"
)

func InitialTokenBalance() sdk.Int {
	return sdk.NewInt(types.Denominator)
}

// Uint64ToBytes converts uint64 to bytes using BigEndian endianness
func Uint64ToBytes(i uint64) []byte {
	b := make([]byte, binary.Size(i))
	binary.BigEndian.PutUint64(b, i)
	return b
}

// BytesToUint64 convert BigEndian unix time from bytes to time.Time
func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// TokenToFloat64 converts token to its float64 representation
func TokenToFloat64(token sdk.Int) float64 {
	if token.IsNil() {
		return 0
	}

	return float64(token.Int64()) / float64(types.Denominator)
}
