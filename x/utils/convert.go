package utils

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/token/types"
)

func InitialTokenBalance() sdk.Int {
	return types.Denominator
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

func UDecToDec(amount sdk.Dec) sdk.Dec {
	if amount.IsNil() || amount.IsZero() {
		return amount
	}

	return amount.QuoInt(types.Denominator)
}
