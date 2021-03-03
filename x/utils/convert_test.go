package utils

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestUint64ToBytes(t *testing.T) {
	i := uint64(101)

	b := Uint64ToBytes(i)
	assert.Equal(t, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65}, b)
}

func TestBytesToUint64(t *testing.T) {
	b := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x65}

	assert.EqualValues(t, 101, BytesToUint64(b))
}

func TestTokenToFloat64(t *testing.T) {
	assert.EqualValues(t, 1, TokenToFloat64(sdk.NewInt(1000000)))
	assert.EqualValues(t, 0, TokenToFloat64(sdk.Int{}))
}
