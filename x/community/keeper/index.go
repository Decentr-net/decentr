package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/Decentr-net/decentr/x/community/types"
)

const (
	createdAtIndexBucket  = "created_at_idx" // key: timestamp+uuid, value: keeper key
	popularityIndexBucket = "popularity_idx" // key: (likes-dislikes)+uuid, value: keeper key
)

type Index interface {
	AddPost(p types.Post) error
	DeletePost(p types.Post) error
	UpdateLikes(old, new types.Post) error

	GetPosts(index string, resolver func([]byte) types.Post, c types.Category, from []byte, limit uint32) ([]types.Post, error)
}

type index struct {
	db *bolt.DB
}

func NewIndex(db *bolt.DB) (Index, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		if err := createIndexBucket(createdAtIndexBucket, tx); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", createdAtIndexBucket, err)
		}

		if err := createIndexBucket(popularityIndexBucket, tx); err != nil {
			return fmt.Errorf("failed to create %s bucket: %w", popularityIndexBucket, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return index{db}, nil
}

func (i index) AddPost(p types.Post) error {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		key := getPostKeeperKeyFromPost(p)

		if err := addPostToIndex(tx.Bucket([]byte(createdAtIndexBucket)), p.Category, getCreateAtIndexKey(p), key); err != nil {
			return err
		}

		if err := addPostToIndex(tx.Bucket([]byte(popularityIndexBucket)), p.Category, getPopularityIndexKey(p), key); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to add post to index: %w", err)
	}

	return nil
}

func (i index) DeletePost(p types.Post) error {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		if err := deletePostFromIndex(tx.Bucket([]byte(createdAtIndexBucket)), p.Category, getCreateAtIndexKey(p)); err != nil {
			return err
		}

		if err := deletePostFromIndex(tx.Bucket([]byte(popularityIndexBucket)), p.Category, getPopularityIndexKey(p)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to delete post from index: %w", err)
	}

	return nil
}

func (i index) UpdateLikes(old, new types.Post) error {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(popularityIndexBucket))
		if err := deletePostFromIndex(bucket, old.Category, getPopularityIndexKey(old)); err != nil {
			return err
		}

		if err := addPostToIndex(bucket, new.Category, getPopularityIndexKey(new), getPostKeeperKeyFromPost(new)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to update post likes to index: %w", err)
	}

	return nil
}

func (i index) GetPosts(index string, resolver func([]byte) types.Post, c types.Category, from []byte, limit uint32) ([]types.Post, error) {
	out := make([]types.Post, 0)

	if err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(index)).Bucket(int64ToBytes(int64(c)))

		c := b.Cursor()
		ik, kk := c.Last()

		if from != nil {
			for ; ik != nil && bytes.Compare(from, ik) != 1; ik, kk = c.Prev() {
			}
		}

		if ik == nil {
			return nil
		}

		for i := uint32(0); i < limit && ik != nil; i++ {
			out = append(out, resolver(kk))
			ik, kk = c.Prev()
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return out, nil
}

func createIndexBucket(name string, tx *bolt.Tx) error {
	b, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return err
	}

	for i := types.UndefinedCategory; i <= types.FitnessAndExerciseCategory; i++ {
		if _, err := b.CreateBucketIfNotExists(int64ToBytes(int64(i))); err != nil {
			return fmt.Errorf("failed to create category bucket: %w", err)
		}
	}

	return nil
}

func addPostToIndex(b *bolt.Bucket, c types.Category, indexKey, keeperKey []byte) error {
	if err := b.Bucket(int64ToBytes(int64(types.UndefinedCategory))).Put(indexKey, keeperKey); err != nil {
		return err
	}

	return b.Bucket(int64ToBytes(int64(c))).Put(indexKey, keeperKey)
}

func deletePostFromIndex(b *bolt.Bucket, c types.Category, indexKey []byte) error {
	if err := b.Bucket(int64ToBytes(int64(types.UndefinedCategory))).Delete(indexKey); err != nil {
		return err
	}

	return b.Bucket(int64ToBytes(int64(c))).Delete(indexKey)
}

func getCreateAtIndexKey(p types.Post) []byte {
	return append(int64ToBytes(p.CreatedAt), p.UUID.Bytes()...)
}

func getPopularityIndexKey(p types.Post) []byte {
	diff := int64(p.LikesCount)
	return append(int64ToBytes(diff), p.UUID.Bytes()...)
}

func int64ToBytes(i int64) []byte {
	b := bytes.NewBuffer(make([]byte, 8))

	_ = binary.Write(b, binary.BigEndian, i) // use BE endianness to have proper sort

	return b.Bytes()
}
