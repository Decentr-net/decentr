package community

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Decentr-net/decentr/x/community/keeper"
	"github.com/Decentr-net/decentr/x/community/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if k.SetParams(ctx, types.DefaultParams()); genState.Params != nil {
		k.SetParams(ctx, *genState.Params)
	}

	for _, v := range genState.Posts {
		k.SetPost(ctx, v)
	}

	for _, v := range genState.Likes {
		k.SetLike(ctx, v)
	}

	for who, addressList := range genState.Following {
		owner, err := sdk.AccAddressFromBech32(who)
		if err != nil {
			panic(fmt.Sprintf("invalid owner: %s", err.Error()))
		}
		for _, whom := range addressList.Address {
			k.Follow(ctx, owner, whom)
		}
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	var posts []types.Post
	k.IteratePosts(ctx, func(p types.Post) (stop bool) {
		posts = append(posts, p)
		return false
	})

	var likes []types.Like
	k.IterateLikes(ctx, func(l types.Like) (stop bool) {
		likes = append(likes, l)
		return false
	})

	following := map[string]types.GenesisState_AddressList{}
	k.IterateFollowings(ctx, func(who, whom sdk.AccAddress) (stop bool) {
		l := following[who.String()]
		l.Address = append(l.Address, whom)
		following[who.String()] = l
		return false
	})
	if len(following) == 0 {
		following = nil
	}

	return &types.GenesisState{
		Params:    &params,
		Posts:     posts,
		Likes:     likes,
		Following: following,
	}
}
