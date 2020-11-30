package keeper

import (
	"testing"

	"github.com/Decentr-net/decentr/x/community/types"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePostLikesCounters(t *testing.T) {
	tt := []struct {
		name           string
		oldWeight      types.LikeWeight
		newWeight      types.LikeWeight
		likesCount     uint32
		dislikeCount   uint32
		wantLikesCount uint32
		wantDikesCount uint32
	}{
		{"zero", types.LikeWeightZero, types.LikeWeightZero, 1, 2, 1, 2},
		{"up=>down", types.LikeWeightUp, types.LikeWeightDown, 1, 2, 0, 3},
		{"down=>up", types.LikeWeightDown, types.LikeWeightUp, 1, 2, 2, 1},
		{"up=>up", types.LikeWeightUp, types.LikeWeightUp, 1, 2, 1, 2},
		{"down=>down", types.LikeWeightDown, types.LikeWeightDown, 1, 2, 1, 2},
		{"zero=>up", types.LikeWeightZero, types.LikeWeightUp, 1, 2, 2, 2},
		{"zero=>down", types.LikeWeightZero, types.LikeWeightDown, 1, 2, 1, 3},
		{"up=>zero", types.LikeWeightUp, types.LikeWeightZero, 1, 2, 0, 2},
		{"down=>zero", types.LikeWeightDown, types.LikeWeightZero, 1, 2, 1, 1},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var post types.Post
			post.LikesCount = tc.likesCount
			post.DislikesCount = tc.dislikeCount

			updatePostLikesCounters(&post, tc.oldWeight, tc.newWeight)

			assert.Equal(t, tc.wantLikesCount, post.LikesCount)
			assert.Equal(t, tc.wantDikesCount, post.DislikesCount)
		})
	}
}
