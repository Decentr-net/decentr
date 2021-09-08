package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type RewardDistribution struct {
	Height int64     `json:"height"`
	Coins  sdk.Coins `json:"coins"`
}
