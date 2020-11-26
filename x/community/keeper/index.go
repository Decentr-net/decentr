package keeper

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/types"
	"github.com/Decentr-net/decentr/x/utils"
)

const (
	createdAtIndexBucket  = "created_at_idx" // key: timestamp+uuid, value: keeper key
	popularityIndexBucket = "popularity_idx" // key: (likes-dislikes)+uuid, value: keeper key
)

type Interval uint8

const (
	InvalidInterval Interval = iota
	DayInterval
	WeekInterval
	MonthInterval
)

var intervals = map[Interval]time.Duration{
	DayInterval:   time.Hour * 24,
	WeekInterval:  time.Hour * 24 * 7,
	MonthInterval: time.Hour * 24 * 31,
}

type Index interface {
	AddPost(p types.Post) error
	DeletePost(p types.Post) error
	UpdateLikes(old, new types.Post) error

	GetRecentPosts(resolver func([]byte) types.Post, c types.Category, from []byte, limit uint32) ([]types.Post, error)
	GetPopularPosts(resolver func([]byte) types.Post, i Interval, c types.Category, from []byte, limit uint32) ([]types.Post, error)

	RemoveUnnecessaryPosts(ctx sdk.Context, t uint64, resolver func([]byte) types.Post)
}

// This index consists of 2 main buckets: created_at_idx and popularity_idx
// Every main bucket consists of categories buckets:
// -- UndefinedCategory(0) is used for global searching
// -- Others buckets are used to search in certain category
// Categories buckets in popularity_idx buckets also consist of interval buckets which used for searching in certain intervals: 1(daily), 2(weekly), 3(monthly)

type index struct {
	db *bolt.DB
}

func NewIndex(db *bolt.DB) (Index, error) {
	if err := db.Update(func(tx *bolt.Tx) error {
		buckets := make([][][]byte, 0)

		for c := types.UndefinedCategory; c <= types.FitnessAndExerciseCategory; c++ {
			buckets = append(buckets, [][]byte{
				[]byte(createdAtIndexBucket),
				utils.Uint64ToBytes(uint64(c)),
			})

			for i := range intervals {
				buckets = append(buckets, [][]byte{
					[]byte(popularityIndexBucket),
					utils.Uint64ToBytes(uint64(c)),
					utils.Uint64ToBytes(uint64(i)),
				})
			}
		}

		for _, b := range buckets {
			if err := createBucket(tx, b); err != nil {
				return fmt.Errorf("failed to create %s bucket: %w", fmt.Sprintln(b), err)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return index{db}, nil
}

func (i index) AddPost(p types.Post) error {
	if err := i.db.Batch(func(tx *bolt.Tx) error {
		key := getPostKeeperKeyFromPost(p)

		for _, b := range getCreatedAtIndexBuckets(p.Category) {
			if err := addPostToIndex(tx, b, getCreateAtIndexKey(p), key); err != nil {
				return err
			}
		}

		for _, b := range getPopularityIndexBuckets(p.Category, p.CreatedAt) {
			if err := addPostToIndex(tx, b, getPopularityIndexKey(p), key); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to add post to index: %w", err)
	}

	return nil
}

func (i index) DeletePost(p types.Post) error {
	if err := i.db.Batch(func(tx *bolt.Tx) error {
		for _, b := range getCreatedAtIndexBuckets(p.Category) {
			if err := deletePostFromIndex(tx, b, getCreateAtIndexKey(p)); err != nil {
				return err
			}
		}

		for _, b := range getPopularityIndexBuckets(p.Category, 0) {
			if err := deletePostFromIndex(tx, b, getPopularityIndexKey(p)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to delete post from index: %w", err)
	}

	return nil
}

func (i index) UpdateLikes(old, new types.Post) error {
	if err := i.db.Batch(func(tx *bolt.Tx) error {
		for _, b := range getPopularityIndexBuckets(old.Category, 0) {
			if err := deletePostFromIndex(tx, b, getPopularityIndexKey(old)); err != nil {
				return err
			}
		}

		for _, b := range getPopularityIndexBuckets(new.Category, new.CreatedAt) {
			if err := addPostToIndex(tx, b, getPopularityIndexKey(new), getPostKeeperKeyFromPost(new)); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to update post likes to index: %w", err)
	}

	return nil
}

func (i index) GetRecentPosts(resolver func([]byte) types.Post, c types.Category, from []byte, limit uint32) ([]types.Post, error) {
	path := [][]byte{
		[]byte(createdAtIndexBucket),
		utils.Uint64ToBytes(uint64(c)),
	}

	return i.getPosts(path, resolver, from, limit)
}

// RemoveUnnecessaryPosts removes posts from inappropriates buckets.
// For example it removes post from daily interval bucket if it's created few days ago
func (i index) RemoveUnnecessaryPosts(ctx sdk.Context, t uint64, resolver func([]byte) types.Post) {
	wg := sync.WaitGroup{}

	flush := func(c types.Category, interval Interval, limit time.Duration) {
		defer wg.Done()

		l := t - uint64(limit/time.Second)
		if err := i.db.Batch(func(tx *bolt.Tx) error {
			c := getBucket(tx, [][]byte{
				[]byte(popularityIndexBucket),
				utils.Uint64ToBytes(uint64(c)),
				utils.Uint64ToBytes(uint64(interval)),
			}).Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				post := resolver(v)

				if post.CreatedAt < l {
					if err := c.Delete(); err != nil {
						ctx.Logger().Error("failed to flush post from index",
							"err", err.Error(),
							"post", fmt.Sprintf("%s/%s", post.Owner, post.UUID),
							"interval", interval,
						)
					}
				}
			}

			return nil
		}); err != nil {
			ctx.Logger().Error("failed to flush index because of executing batch failure",
				"err", err.Error(),
				"interval", interval,
			)
		}
	}

	for i := types.UndefinedCategory; i <= types.FitnessAndExerciseCategory; i++ {
		for b, d := range intervals {
			wg.Add(1)
			go flush(i, b, d)
		}
	}

	wg.Wait()
}

func (i index) GetPopularPosts(resolver func([]byte) types.Post, interval Interval, c types.Category, from []byte, limit uint32) ([]types.Post, error) {
	path := [][]byte{
		[]byte(popularityIndexBucket),
		utils.Uint64ToBytes(uint64(c)),
		utils.Uint64ToBytes(uint64(interval)),
	}

	return i.getPosts(path, resolver, from, limit)
}

func (i index) getPosts(path [][]byte, resolver func([]byte) types.Post, from []byte, limit uint32) ([]types.Post, error) {
	out := make([]types.Post, 0)

	if err := i.db.View(func(tx *bolt.Tx) error {
		c := getBucket(tx, path).Cursor()
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

func createBucket(tx *bolt.Tx, p [][]byte) error {
	if len(p) == 0 {
		panic("path is empty")
	}

	b, err := tx.CreateBucketIfNotExists(p[0])
	if err != nil {
		return err
	}

	for _, v := range p[1:] {
		if b, err = b.CreateBucketIfNotExists(v); err != nil {
			return err
		}
	}

	return nil
}

func getCreatedAtIndexBuckets(c types.Category) [][][]byte {
	return [][][]byte{
		{[]byte(createdAtIndexBucket), utils.Uint64ToBytes(uint64(types.UndefinedCategory))},
		{[]byte(createdAtIndexBucket), utils.Uint64ToBytes(uint64(c))},
	}
}

func getPopularityIndexBuckets(c types.Category, createdAt uint64) [][][]byte {
	p := [][][]byte{
		{[]byte(popularityIndexBucket), utils.Uint64ToBytes(uint64(types.UndefinedCategory))},
		{[]byte(popularityIndexBucket), utils.Uint64ToBytes(uint64(c))},
	}

	out := make([][][]byte, 0, len(intervals)*len(p))

	for _, b := range p {
		for i, v := range intervals {
			if createdAt == 0 ||
				createdAt+uint64(v/time.Second) > uint64(time.Now().Unix()) {
				out = append(out, append(b, utils.Uint64ToBytes(uint64(i))))
			}
		}
	}

	return out
}

func addPostToIndex(tx *bolt.Tx, path [][]byte, indexKey, keeperKey []byte) error {
	if err := getBucket(tx, path).Put(indexKey, keeperKey); err != nil {
		return err
	}

	return nil
}

func deletePostFromIndex(tx *bolt.Tx, path [][]byte, indexKey []byte) error {
	if err := getBucket(tx, path).Delete(indexKey); err != nil {
		return err
	}

	return nil
}

func getCreateAtIndexKey(p types.Post) []byte {
	return append(utils.Uint64ToBytes(p.CreatedAt), p.UUID.Bytes()...)
}

func getPopularityIndexKey(p types.Post) []byte {
	return append(utils.Uint64ToBytes(uint64(p.LikesCount)), p.UUID.Bytes()...)
}

func getBucket(tx *bolt.Tx, path [][]byte) *bolt.Bucket {
	if len(path) == 0 {
		panic("path is empty")
	}

	b := tx.Bucket(path[0])
	for _, v := range path[1:] {
		b = b.Bucket(v)
	}

	return b
}
