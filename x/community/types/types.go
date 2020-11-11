package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"
)

type Post struct {
	UUID          uuid.UUID      `json:"uuid"`
	Owner         sdk.AccAddress `json:"owner"`
	Title         string         `json:"title"`
	PreviewImage  string         `json:"previewImage"`
	Text          string         `json:"text"`
	LikesCount    uint32         `json:"likesCount"`
	DislikesCount uint32         `json:"dislikesCount"`
	CreatedAt     time.Time      `json:"createdAt"`
}
