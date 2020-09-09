package keeper

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/pdv/types"
)

const (
	ownersBucket = "owners"
)

type Stats interface {
	AddPDV(pdv types.PDV, token sdk.Int)
	ListPDV(owner sdk.AccAddress, from *time.Time, limit uint) ([]types.PDV, error)
	GetStats(owner sdk.AccAddress) (map[time.Time]sdk.Int, error)
}

type stats struct {
	db *bolt.DB
}

type statsItem struct {
	Address string        `json:"address"`
	Type    types.PDVType `json:"type"`
	Token   sdk.Int       `json:"token"`
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

	return &stats{db}, nil
}

func (s *stats) AddPDV(pdv types.PDV, token sdk.Int) {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket))
		b, err := b.CreateBucketIfNotExists(pdv.Owner)
		if err != nil {
			return fmt.Errorf("failed to create owner bucket: %w", err)
		}

		v, err := json.Marshal(statsItem{
			Address: pdv.Address,
			Type:    pdv.Type,
			Token:   token,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal stats item: %w", err)
		}

		if err := b.Put([]byte(pdv.Timestamp.Format(time.RFC3339)), v); err != nil {
			return fmt.Errorf("failed to put stats item: %w", err)
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("failed to add pdv to stats: %w", err))
	}
}

func (s *stats) ListPDV(owner sdk.AccAddress, from *time.Time, limit uint) ([]types.PDV, error) {
	res := make([]types.PDV, 0, limit)
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket)).Bucket(owner)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		k, v := c.Last()
		if from != nil {
			k, _ = c.Seek([]byte(from.Format(time.RFC3339)))
			if k == nil {
				k, v = c.Last()
			} else {
				k, v = c.Prev()
			}
		}

		var i uint

		for ; k != nil && i < limit; k, v = c.Prev() {
			timestamp, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return fmt.Errorf("invalid stats item key: %s", k)
			}

			var si statsItem
			if err := json.Unmarshal(v, &si); err != nil {
				return fmt.Errorf("failed to unmarshal stats item: %s", k)
			}

			res = append(res, types.PDV{
				Timestamp: timestamp,
				Address:   si.Address,
				Owner:     owner,
				Type:      si.Type,
			})
			i++
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *stats) GetStats(owner sdk.AccAddress) (map[time.Time]sdk.Int, error) {
	res := make(map[time.Time]sdk.Int, 365)
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket)).Bucket(owner)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		d := time.Time{}
		t := sdk.NewInt(0)

		for k, v := c.First(); k != nil; k, v = c.Next() {
			timestamp, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return fmt.Errorf("invalid stats item key: %s", k)
			}

			if timestamp.Truncate(time.Hour*24) != d {
				res[d] = t
				d = timestamp.Truncate(time.Hour * 24)
			}

			var si statsItem
			if err := json.Unmarshal(v, &si); err != nil {
				return fmt.Errorf("failed to unmarshal stats item: %s", k)
			}

			t = t.Add(si.Token)
		}

		res[d] = t

		return nil
	}); err != nil {
		return nil, err
	}

	return res, nil
}
