package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/pdv/types"
	"github.com/Decentr-net/decentr/x/utils"
)

const (
	ownersBucket = "owners"
)

type Index interface {
	AddPDV(pdv types.PDV) error
	ListPDV(owner sdk.AccAddress, from *uint64, limit uint) ([]types.PDV, error)
}

type index struct {
	db *bolt.DB
}

type statsItem struct {
	Address string        `json:"address"`
	Type    types.PDVType `json:"type"`
}

func NewIndex(db *bolt.DB) (Index, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(ownersBucket)); err != nil {
			return fmt.Errorf("failed to create owners bucket: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &index{db}, nil
}

func (s *index) AddPDV(pdv types.PDV) error {
	if err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket))
		b, err := b.CreateBucketIfNotExists(pdv.Owner)
		if err != nil {
			return fmt.Errorf("failed to create owner bucket: %w", err)
		}

		v, err := json.Marshal(statsItem{
			Address: pdv.Address,
			Type:    pdv.Type,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal index item: %w", err)
		}

		if err := b.Put(utils.Uint64ToBytes(pdv.Timestamp), v); err != nil {
			return fmt.Errorf("failed to put index item: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to add pdv to index: %w", err)
	}

	return nil
}

func (s *index) ListPDV(owner sdk.AccAddress, from *uint64, limit uint) ([]types.PDV, error) {
	res := make([]types.PDV, 0, limit)
	if err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ownersBucket)).Bucket(owner)
		if b == nil {
			return nil
		}

		c := b.Cursor()
		k, v := c.Last()
		if from != nil {
			k, _ = c.Seek(utils.Uint64ToBytes(*from))
			if k == nil {
				k, v = c.Last()
			} else {
				k, v = c.Prev()
			}
		}

		var i uint

		for ; k != nil && i < limit; k, v = c.Prev() {
			var si statsItem
			if err := json.Unmarshal(v, &si); err != nil {
				return fmt.Errorf("failed to unmarshal index item: %s", k)
			}

			res = append(res, types.PDV{
				Timestamp: utils.BytesToUint64(k),
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
