package keeper

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

func getIndex() Index {
	d, err := ioutil.TempDir("", "*")
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(fmt.Sprintf("%s/file.db", d), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	i, err := NewIndex(db)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

func TestStats_AddPDV(t *testing.T) {
	i := getIndex()

	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 7}
	pdv := types.PDV{
		Timestamp: 1578009600,
		Address:   "address",
		Owner:     owner,
		Type:      types.PDVTypeCookie,
	}
	require.NoError(t, i.AddPDV(pdv))

	p, err := i.ListPDV(sdk.AccAddress{1, 2, 3, 4, 5, 6, 7}, nil, 5)
	require.NoError(t, err)
	require.Len(t, p, 1)
	assert.Equal(t, pdv, p[0])
}

func TestStats_ListPDV(t *testing.T) {
	i := getIndex()

	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 1}
	for j := 0; j < 20; j++ {
		require.NoError(t, i.AddPDV(types.PDV{
			Timestamp: 1578009600 + uint64(j),
			Address:   "address",
			Owner:     owner,
			Type:      types.PDVTypeCookie,
		}))
	}

	p, err := i.ListPDV(owner, nil, 30)
	require.NoError(t, err)
	assert.Len(t, p, 20)

	p, err = i.ListPDV(owner, nil, 10)
	require.NoError(t, err)
	assert.Len(t, p, 10)
	assert.EqualValues(t, 1578009619, p[0].Timestamp)

	p, err = i.ListPDV(owner, &p[9].Timestamp, 10)
	require.NoError(t, err)
	assert.Len(t, p, 10)
	assert.EqualValues(t, 1578009609, p[0].Timestamp)
}
