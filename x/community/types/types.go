package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"
)

type Category uint8

const (
	UndefinedCategory Category = iota
	WorldNewsCategory
	TravelAndTourismCategory
	ScienceAndTechnologyCategory
	StrangeWorldCategory
	ArtsAndEntertainmentCategory
	WritersAndWritingCategory
	HealthAndFitnessCategory
	CryptoAndBlockchainCategory
	SportsCategory
)

type LikeWeight int8

const (
	LikeWeightUp   LikeWeight = 1
	LikeWeightZero LikeWeight = 0
	LikeWeightDown LikeWeight = -1
)

type Post struct {
	UUID          uuid.UUID      `json:"uuid"`
	Owner         sdk.AccAddress `json:"owner"`
	Title         string         `json:"title"`
	PreviewImage  string         `json:"previewImage"`
	Category      Category       `json:"category"`
	Text          string         `json:"text"`
	LikesCount    uint32         `json:"likesCount"`
	DislikesCount uint32         `json:"dislikesCount"`
	PDV           sdk.Int        `json:"pdv"`
	CreatedAt     uint64         `json:"createdAt"`
}

type Like struct {
	Owner     sdk.AccAddress `json:"owner"`
	PostOwner sdk.AccAddress `json:"postOwner"`
	PostUUID  uuid.UUID      `json:"postUuid"`
	Weight    LikeWeight     `json:"weight"`
}
