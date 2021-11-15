package v150

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CommunityState struct {
	Params struct {
		Moderators []string `json:"moderators" yaml:"moderators"`
		FixedGas   struct {
			CreatePost sdk.Gas `json:"create_post" yaml:"create_post"`
			DeletePost sdk.Gas `json:"delete_post" yaml:"delete_post"`
			SetLike    sdk.Gas `json:"set_like" yaml:"set_like"`
			Follow     sdk.Gas `json:"follow" yaml:"follow"`
			Unfollow   sdk.Gas `json:"unfollow" yaml:"unfollow"`
		} `json:"fixed_gas" yaml:"fixed_gas"`
	} `json:"params"`
	Posts []struct {
		UUID         string         `json:"uuid"`
		Owner        sdk.AccAddress `json:"owner"`
		Title        string         `json:"title"`
		PreviewImage string         `json:"previewImage"`
		Category     uint8          `json:"category"`
		Text         string         `json:"text"`
	} `json:"posts"`
	Likes []struct {
		Owner     sdk.AccAddress `json:"owner"`
		PostOwner sdk.AccAddress `json:"postOwner"`
		PostUUID  string         `json:"postUuid"`
		Weight    int8           `json:"weight"`
	} `json:"likes"`
	Followers map[string][]string `json:"followers"`
}

type OperationsState struct {
	Params struct {
		Supervisors []string `json:"supervisors" yaml:"supervisors"`
		FixedGas    struct {
			ResetAccount      sdk.Gas `json:"reset_account" yaml:"reset_account"`
			BanAccount        sdk.Gas `json:"ban_account" yaml:"ban_account"`
			DistributeRewards sdk.Gas `json:"distribute_rewards" yaml:"distribute_rewards"`
		} `json:"fixed_gas" yaml:"fixed_gas"`
		MinGasPrice sdk.DecCoin `json:"min_gas_price" yaml:"min_gas_price"`
	} `json:"params"`
}

type TokenState struct {
	Params struct {
		RewardsBlockInterval int64 `json:"rewards_block_interval" yaml:"rewards_block_interval"`
	} `json:"params"`
	Balances map[string]sdk.Int `json:"balances"`
	Deltas   map[string]sdk.Int `json:"deltas"`
	History  map[string][]struct {
		Height int64     `json:"height"`
		Coins  sdk.Coins `json:"coins"`
	} `json:"history"`
}
