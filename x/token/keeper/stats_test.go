package keeper

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Decentr-net/decentr/x/utils"
)

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

func TestStats_AddToken(t *testing.T) {
	s := getStats()

	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 7}
	timestamp := uint64(time.Date(time.Now().Year(), 1, 3, 4, 5, 0, 0, time.UTC).Unix())

	require.NoError(t, s.AddToken(owner, timestamp, sdk.NewIntFromUint64(1)))

	stats, err := s.GetStats(owner)
	require.NoError(t, err)
	for i, v := range stats {
		if i == truncateUnixTime(timestamp, time.Hour*24) {
			assert.EqualValues(t, 0.0000001, v)
		} else {
			assert.EqualValues(t, float64(0), v)
		}
	}
}

func TestStats_GetStats(t *testing.T) {
	s := getStats()
	owner := sdk.AccAddress{1, 2, 3, 4, 5, 6, 200}
	for i := 1; i <= 31; i++ {
		timestamp := uint64(time.Date(time.Now().Year(), 1, i, 4, 5, 0, 0, time.UTC).Unix())
		require.NoError(t, s.AddToken(owner, timestamp, sdk.NewIntFromUint64(uint64(i))))
	}

	sum := 0
	res, err := s.GetStats(owner)
	require.NoError(t, err)
	require.Len(t, res, 32)
	assert.EqualValues(t, 0, res[0])
	for i := 1; i <= 31; i++ {
		tm := uint64(time.Date(2020, 1, i, 0, 0, 0, 0, time.UTC).Unix())

		sum += i

		assert.EqualValues(t, utils.TokenToFloat64(sdk.NewInt(int64(sum))), res[tm])
	}
}
