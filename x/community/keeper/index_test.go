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

	types2 "github.com/Decentr-net/decentr/x/community/types"
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

type localResolver map[string]types2.Post

func (r localResolver) resolve(k []byte) types2.Post {
	return r[string(k)]
}

func (r localResolver) add(i Index, p types2.Post) error {
	r[string(getPostKeeperKey(p.Owner, p.UUID))] = p
	return i.AddPost(p)
}

func (r localResolver) updateLikes(i Index, old, new types2.Post) error {
	delete(r, string(getPostKeeperKey(old.Owner, new.UUID)))
	r[string(getPostKeeperKey(new.Owner, new.UUID))] = new

	return i.UpdateLikes(old, new)
}

func TestIndex_AddPost(t *testing.T) {
	i := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	require.NoError(t, r.add(i, types2.Post{
		Owner:      testOwner,
		Category:   types2.FitnessAndExerciseCategory,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp,
		LikesCount: 2,
	}))
	require.NoError(t, r.add(i, types2.Post{
		Owner:      testOwner,
		Category:   types2.HealthAndCultureCategory,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp + 1,
		LikesCount: 1,
	}))

	p, err := i.GetRecentPosts(r.resolve, types2.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 2)
	require.EqualValues(t, timestamp+1, p[0].CreatedAt)
	require.EqualValues(t, timestamp, p[1].CreatedAt)

	p, err = i.GetPopularPosts(r.resolve, MonthInterval, types2.HealthAndCultureCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 1)
	require.EqualValues(t, timestamp+1, p[0].CreatedAt)
}

func TestIndex_DeletePost(t *testing.T) {
	i := getIndex()
	r := make(localResolver)

	p := types2.Post{
		Owner:      testOwner,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  uint64(time.Now().Unix()),
		LikesCount: 2,
	}
	require.NoError(t, r.add(i, p))

	require.NoError(t, i.DeletePost(p))

	posts, err := i.GetPopularPosts(r.resolve, MonthInterval, types2.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Empty(t, posts)
}

func TestIndex_UpdateLikes(t *testing.T) {
	i := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	old := types2.Post{
		Owner:      testOwner,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp,
		LikesCount: 2,
	}
	require.NoError(t, r.add(i, old))

	require.NoError(t, r.add(i, types2.Post{
		Owner:      testOwner,
		UUID:       uuid.Must(uuid.NewV1()),
		CreatedAt:  timestamp + 1,
		LikesCount: 5,
	}))

	p, err := i.GetPopularPosts(r.resolve, MonthInterval, types2.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 2)
	require.EqualValues(t, timestamp+1, p[0].CreatedAt)
	require.EqualValues(t, timestamp, p[1].CreatedAt)

	new := old
	new.LikesCount = 20
	require.NoError(t, r.updateLikes(i, old, new))

	p, err = i.GetPopularPosts(r.resolve, MonthInterval, types2.UndefinedCategory, nil, 10)
	require.NoError(t, err)
	require.Len(t, p, 2)
	require.EqualValues(t, timestamp, p[0].CreatedAt)
	require.EqualValues(t, timestamp+1, p[1].CreatedAt)
}

func TestIndex_getPosts(t *testing.T) {
	i := getIndex()

	r := make(localResolver)

	p := types2.Post{
		Owner:     testOwner,
		UUID:      uuid.Must(uuid.NewV1()),
		CreatedAt: 4,
	}
	require.NoError(t, r.add(i, p))
	require.NoError(t, r.add(i, types2.Post{
		Owner:     testOwner,
		UUID:      uuid.Must(uuid.NewV1()),
		CreatedAt: 3,
	}))

	posts, err := i.GetRecentPosts(r.resolve, types2.UndefinedCategory, getCreateAtIndexKey(p), 10)
	require.NoError(t, err)
	require.Len(t, posts, 1)
	require.EqualValues(t, 3, posts[0].CreatedAt)
}

func TestIndex_Flush(t *testing.T) {
	t.Parallel()

	i := getIndex()
	r := make(localResolver)

	timestamp := uint64(time.Now().Unix())

	for c := types2.WorldNewsCategory; c <= types2.FitnessAndExerciseCategory; c++ {
		for j := 0; j < 10; j++ {
			require.NoError(t, r.add(i, types2.Post{
				Owner:      testOwner,
				Category:   c,
				UUID:       uuid.Must(uuid.NewV1()),
				CreatedAt:  timestamp + uint64(j*267800),
				LikesCount: uint32(j),
			}))
		}
		require.NoError(t, r.add(i, types2.Post{
			Owner:      testOwner,
			Category:   c,
			UUID:       uuid.Must(uuid.NewV1()),
			CreatedAt:  timestamp + 2678500, // add month with tail
			LikesCount: 1,
		}))
	}

	i.RemoveUnnecessaryPosts(sdk.Context{}, timestamp+2678400*2, r.resolve) // we will pretend that 2 month have passed

	p, err := i.GetPopularPosts(r.resolve, MonthInterval, types2.UndefinedCategory, nil, 20)
	require.NoError(t, err)
	require.Len(t, p, 6)
}
