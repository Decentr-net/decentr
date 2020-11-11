package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"

	"github.com/Decentr-net/decentr/x/community/types"
)

const (
	createdAtIndexBucket  = "created_at_idx" // key: timestamp+uuid, value: keeper key
	popularityIndexBucket = "popularity_idx" // key: (likes-dislikes)+uuid, value: keeper key
)

type Index interface {
	AddPost(p types.Post)
	DeletePost(p types.Post)
	UpdateLikes(old, new types.Post)
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

	return &index{db}, nil
}

func (i index) AddPost(p types.Post) {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		key := getPostKeeperKey(p)

		if err := addPostToIndex(tx.Bucket([]byte(createdAtIndexBucket)), p.Category, getCreateAtIndexKey(p), key); err != nil {
			return err
		}

		if err := addPostToIndex(tx.Bucket([]byte(popularityIndexBucket)), p.Category, getPopularityIndexKey(p), key); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("failed to add post to index: %w", err))
	}
}

func (i index) DeletePost(p types.Post) {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		if err := deletePostFromIndex(tx.Bucket([]byte(createdAtIndexBucket)), p.Category, getCreateAtIndexKey(p)); err != nil {
			return err
		}

		if err := deletePostFromIndex(tx.Bucket([]byte(popularityIndexBucket)), p.Category, getPopularityIndexKey(p)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("failed to delete post from index: %w", err))
	}
}

func (i index) UpdateLikes(old, new types.Post) {
	if err := i.db.Update(func(tx *bolt.Tx) error {
		if err := deletePostFromIndex(tx.Bucket([]byte(popularityIndexBucket)), old.Category, getPopularityIndexKey(old)); err != nil {
			return err
		}

		if err := addPostToIndex(tx.Bucket([]byte(popularityIndexBucket)), new.Category, getPopularityIndexKey(new), getPostKeeperKey(new)); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(fmt.Errorf("failed to update post likes to index: %w", err))
	}
}

func createIndexBucket(name string, tx *bolt.Tx) error {
	b, err := tx.CreateBucketIfNotExists([]byte(name))
	if err != nil {
		return err
	}

	for i := types.UndefinedCategory; i <= types.FitnessAndExerciseCategory; i++ {
		if _, err := b.CreateBucketIfNotExists([]byte(strconv.Itoa(int(i)))); err != nil {
			return fmt.Errorf("failed to create category bucket: %w", err)
		}
	}

	return nil
}

func addPostToIndex(b *bolt.Bucket, c types.Category, indexKey, keeperKey []byte) error {
	if err := b.Bucket([]byte(strconv.Itoa(int(types.UndefinedCategory)))).Put(indexKey, keeperKey); err != nil {
		return err
	}

	return b.Bucket([]byte(strconv.Itoa(int(c)))).Put(indexKey, keeperKey)
}

func deletePostFromIndex(b *bolt.Bucket, c types.Category, indexKey []byte) error {
	if err := b.Bucket([]byte(strconv.Itoa(int(types.UndefinedCategory)))).Delete(indexKey); err != nil {
		return err
	}

	return b.Bucket([]byte(strconv.Itoa(int(c)))).Delete(indexKey)
}

func getCreateAtIndexKey(p types.Post) []byte {
	return append(int64ToBytes(p.CreatedAt), p.UUID.Bytes()...)
}

func getPopularityIndexKey(p types.Post) []byte {
	diff := int64(p.LikesCount) - int64(p.DislikesCount)
	return append(int64ToBytes(diff), p.UUID.Bytes()...)
}

func int64ToBytes(i int64) []byte {
	b := bytes.NewBuffer(make([]byte, 8))

	_ = binary.Write(b, binary.BigEndian, i) // use BE endianness to have proper sort

	return b.Bytes()
}
