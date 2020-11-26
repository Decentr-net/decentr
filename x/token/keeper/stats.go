package keeper

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/utils"
)

const ownersBucket = "owners"

type Stats interface {
	AddToken(owner sdk.AccAddress, timestamp uint64, token sdk.Int) error
	GetStats(owner sdk.AccAddress) (map[uint64]float64, error)
}

type stats struct {
	cdc *codec.Codec
	db  *bolt.DB
}

func NewStats(db *bolt.DB) (Stats, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(ownersBucket)); err != nil {
			return fmt.Errorf("failed to create owners bucket: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &stats{
		cdc: codec.New(),
		db:  db,
	}, nil
}

func (s *stats) AddToken(owner sdk.AccAddress, timestamp uint64, amount sdk.Int) error {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket))
		b, err := b.CreateBucketIfNotExists(owner)
		if err != nil {
			return fmt.Errorf("failed to create owner bucket: %w", err)
		}

		v, err := s.cdc.MarshalBinaryBare(amount)
		if err != nil {
			return fmt.Errorf("failed to marshall amount: %w", err)
		}
		if err := b.Put(utils.Uint64ToBytes(timestamp), v); err != nil {
			return fmt.Errorf("failed to put stats item: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to add pdv to stats: %w", err)
	}

	return nil
}

func (s *stats) GetStats(owner sdk.AccAddress) (map[uint64]float64, error) {
	res := make(map[uint64]float64, 365)
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket)).Bucket(owner)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		d := uint64(0)
		t := sdk.NewInt(0)

		for k, v := c.First(); k != nil; k, v = c.Next() {
			timestamp := utils.BytesToUint64(k)

			truncated := truncateUnixTime(timestamp, time.Hour*24)

			if d == 0 {
				d = truncated
			}

			if truncated != d {
				res[d] = utils.TokenToFloat64(t)
				d = truncated
			}

			var amount sdk.Int
			if err := s.cdc.UnmarshalBinaryBare(v, &amount); err != nil {
				return fmt.Errorf("failed to unmarshal stats item for %s: %w", hex.EncodeToString(k), err)
			}

			t = t.Add(amount)
		}

		res[d] = utils.TokenToFloat64(t)

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

func truncateUnixTime(t uint64, d time.Duration) uint64 {
	return t - t%(uint64)(d/time.Second)
}
