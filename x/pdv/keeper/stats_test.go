package keeper

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	tokenkeeper "github.com/Decentr-net/decentr/x/token/keeper"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

var s = getStats()

func getStats() Stats {
	d, err := ioutil.TempDir("", "*")
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(fmt.Sprintf("%s/file.db", d), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	stats, err := NewStats(db)
	if err != nil {
		log.Fatal(err)
	}

	return stats
}

func TestStats_AddPDV(t *testing.T) {
	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 7}
	pdv := types.PDV{
		Timestamp: time.Date(2020, 2, 3, 4, 5, 6, 0, time.UTC),
		Address:   "address",
		Owner:     owner,
		Type:      types.PDVTypeCookie,
	}
	s.AddPDV(pdv, sdk.NewIntFromUint64(1))

	p, err := s.ListPDV(sdk.AccAddress{1, 2, 3, 4, 5, 6, 7}, nil, 5)
	require.NoError(t, err)
	require.Len(t, p, 1)
	assert.Equal(t, pdv, p[0])
}

func TestStats_ListPDV(t *testing.T) {
	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 1}
	for i := 0; i < 20; i++ {
		s.AddPDV(types.PDV{
			Timestamp: time.Date(2020, 2, 3, 4, 5, i, 0, time.UTC),
			Address:   "address",
			Owner:     owner,
			Type:      types.PDVTypeCookie,
		}, sdk.NewIntFromUint64(uint64(i)))
	}

	p, err := s.ListPDV(owner, nil, 30)
	require.NoError(t, err)
	assert.Len(t, p, 20)

	p, err = s.ListPDV(owner, nil, 10)
	require.NoError(t, err)
	assert.Len(t, p, 10)
	assert.Equal(t, time.Date(2020, 2, 3, 4, 5, 19, 0, time.UTC), p[0].Timestamp)

	p, err = s.ListPDV(owner, &p[9].Timestamp, 10)
	require.NoError(t, err)
	assert.Len(t, p, 10)
	assert.Equal(t, time.Date(2020, 2, 3, 4, 5, 9, 0, time.UTC), p[0].Timestamp)
}

func TestStats_GetStats(t *testing.T) {
	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 200}
	for i := 1; i <= 31; i++ {
		s.AddPDV(types.PDV{
			Timestamp: time.Date(2020, 1, i, 4, 5, 0, 0, time.UTC),
			Address:   "address",
			Owner:     owner,
			Type:      types.PDVTypeCookie,
		}, sdk.NewIntFromUint64(uint64(i)))
	}

	sum := 0
	res, err := s.GetStats(owner)
	require.NoError(t, err)
	require.Len(t, res, 32)
	assert.EqualValues(t, 0, res[time.Time{}])
	for i := 1; i <= 31; i++ {
		tm := time.Date(2020, 1, i, 0, 0, 0, 0, time.UTC)

		sum += i

		assert.EqualValues(t, tokenkeeper.TokenToFloat64(sdk.NewInt(int64(sum))), res[tm])
	}
}
