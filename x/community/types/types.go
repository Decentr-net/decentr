package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gofrs/uuid"
)

type Post struct {
	UUID         uuid.UUID      `json:"uuid"`
	Owner        sdk.AccAddress `json:"owner"`
	Title        string         `json:"title"`
	PreviewImage string         `json:"preview_image"`
	Text         string         `json:"text"`
	Tags         []string       `json:"tags"`
	Likes        int32          `json:"likes"`
	CreatedAt    time.Time      `json:"created_at"`
}
