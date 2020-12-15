package keeper

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"

	"github.com/stretchr/testify/require"

	"github.com/Decentr-net/decentr/x/community/types"
)

var testOwner = sdk.AccAddress{1, 2, 3, 4, 5, 6, 7}

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

type localResolver map[string]types.Post

func (r localResolver) resolve(k []byte) types.Post {
	return r[string(k)]
}

func (r localResolver) add(i Index, p types.Post) error {
	r[string(getPostKeeperKey(p.Owner, p.UUID))] = p
	return i.AddPost(p)
}

func (r localResolver) updateLikes(i Index, old, new types.Post) error {
	delete(r, string(getPostKeeperKey(old.Owner, new.UUID)))
	r[string(getPostKeeperKey(new.Owner, new.UUID))] = new

	return i.UpdateLikes(old, new)
}

func TestIndex_AddPost(t *testing.T) {
	i := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	require.NoError(t, r.add(i, types.Post{
		Owner:      testOwner,
		Category:   types.FitnessAndExerciseCategory,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp,
		LikesCount: 2,
	}))
	require.NoError(t, r.add(i, types.Post{
		Owner:      testOwner,
		Category:   types.HealthAndCultureCategory,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp + 1,
		LikesCount: 1,
	}))

	p, err := i.GetRecentPosts(r.resolve, types.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 2)
	require.EqualValues(t, timestamp+1, p[0].CreatedAt)
	require.EqualValues(t, timestamp, p[1].CreatedAt)

	p, err = i.GetPopularPosts(r.resolve, MonthInterval, types.HealthAndCultureCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 1)
	require.EqualValues(t, timestamp+1, p[0].CreatedAt)
}

func TestIndex_Add10Posts(t *testing.T) {
	const num = 10

	index := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	for i := 0; i < num; i++ {
		require.NoError(t, r.add(index, types.Post{
			Owner:      testOwner,
			Category:   types.FitnessAndExerciseCategory,
			UUID:       uuid.Must(uuid.NewV1()),
			CreatedAt:  timestamp,
			LikesCount: 0,
		}))
	}

	p, err := index.GetRecentPosts(r.resolve, types.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, num)

	p, err = index.GetRecentPosts(r.resolve, types.FitnessAndExerciseCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, num)
}

func TestIndex_DeletePost(t *testing.T) {
	i := getIndex()
	r := make(localResolver)

	p := types.Post{
		Owner:      testOwner,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  uint64(time.Now().Unix()),
		LikesCount: 2,
	}
	require.NoError(t, r.add(i, p))

	require.NoError(t, i.DeletePost(p))

	posts, err := i.GetPopularPosts(r.resolve, MonthInterval, types.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Empty(t, posts)
}

func TestIndex_UpdateLikes(t *testing.T) {
	i := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	old := types.Post{
		Owner:      testOwner,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp,
		LikesCount: 2,
	}
	require.NoError(t, r.add(i, old))

	require.NoError(t, r.add(i, types.Post{
		Owner:      testOwner,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp + 1,
		LikesCount: 5,
	}))

	p, err := i.GetPopularPosts(r.resolve, MonthInterval, types.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 2)
	require.EqualValues(t, timestamp+1, p[0].CreatedAt)
	require.EqualValues(t, timestamp, p[1].CreatedAt)

	new := old
	new.LikesCount = 20
	require.NoError(t, r.updateLikes(i, old, new))

	p, err = i.GetPopularPosts(r.resolve, MonthInterval, types.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 2)
	require.EqualValues(t, timestamp, p[0].CreatedAt)
	require.EqualValues(t, timestamp+1, p[1].CreatedAt)
}

func TestIndex_getPosts(t *testing.T) {
	i := getIndex()

	r := make(localResolver)

	p := types.Post{
		Owner:     testOwner,
		UUID:      uuid.Must(uuid.NewV1()),
		CreatedAt: 4,
	}
	require.NoError(t, r.add(i, p))
	require.NoError(t, r.add(i, types.Post{
		Owner:     testOwner,
		UUID:      uuid.Must(uuid.NewV1()),
		CreatedAt: 3,
	}))

	posts, err := i.GetRecentPosts(r.resolve, types.UndefinedCategory, getCreateAtIndexKey(p), 10)
	require.NoError(t, err)
	require.Len(t, posts, 1)
	require.EqualValues(t, 3, posts[0].CreatedAt)
}

func TestIndex_Flush(t *testing.T) {
	t.Parallel()

	i := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	for c := types.WorldNewsCategory; c <= types.FitnessAndExerciseCategory; c++ {
		for j := 0; j < 10; j++ {
			require.NoError(t, r.add(i, types.Post{
				Owner:      testOwner,
				Category:   c,
				UUID:       uuid.Must(uuid.NewV1()),
				CreatedAt:  timestamp + uint64(j*267800),
				LikesCount: uint32(j),
			}))
		}
		require.NoError(t, r.add(i, types.Post{
			Owner:      testOwner,
			Category:   c,
			UUID:       uuid.Must(uuid.NewV1()),
			CreatedAt:  timestamp + 2678500, // add month with tail
			LikesCount: 1,
		}))
	}

	i.RemoveUnnecessaryPosts(sdk.Context{}, timestamp+2678400*2, r.resolve) // we will pretend that 2 month have passed

	p, err := i.GetPopularPosts(r.resolve, MonthInterval, types.UndefinedCategory, nil, 20)
	require.NoError(t, err)
	require.Len(t, p, 6)
}

func TestIndex_AddLike(t *testing.T) {
	t.Parallel()

	i := getIndex()

	l := types.Like{
		Owner:     testOwner,
		PostOwner: testOwner,
		PostUUID:  uuid.Must(uuid.NewV4()),
		Weight:    types.LikeWeightUp,
	}
	require.NoError(t, i.AddLike(l))

	ll, err := i.GetUserLikedPosts(testOwner)
	require.NoError(t, err)
	require.Len(t, ll, 1)
	require.Equal(t, ll[fmt.Sprintf("%s/%s", l.PostOwner, l.PostUUID)], types.LikeWeightUp)

	l.Weight = types.LikeWeightZero
	require.NoError(t, i.AddLike(l))
	ll, err = i.GetUserLikedPosts(testOwner)
	require.NoError(t, err)
	require.Len(t, ll, 0)
}

func TestIndex_GetUserLikedPosts(t *testing.T) {
	t.Parallel()

	i := getIndex()

	ll, err := i.GetUserLikedPosts(testOwner)
	require.NoError(t, err)
	require.Len(t, ll, 0)

	for j := 0; j < 20; j++ {
		require.NoError(t, i.AddLike(types.Like{
			Owner:     testOwner,
			PostOwner: testOwner,
			PostUUID:  uuid.Must(uuid.NewV4()),
			Weight:    types.LikeWeightUp,
		}))
	}

	ll, err = i.GetUserLikedPosts(testOwner)
	require.NoError(t, err)
	require.Len(t, ll, 20)
}
