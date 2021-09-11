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

func TestUDecToDec(t *testing.T) {
	assert.Equal(t, sdk.NewDec(1), UDecToDec(sdk.NewDec(1000000)))
	assert.Equal(t, sdk.Dec{}, UDecToDec(sdk.Dec{}))
	assert.Equal(t, sdk.NewDec(0), UDecToDec(sdk.NewDec(0)))
}
