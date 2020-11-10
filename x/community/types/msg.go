package types

import (
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gofrs/uuid"
)

const (
	maxTitleLength = 150
	maxPostLength  = 5000
	minPostLength  = 15
	maxTagsCount   = 5
	maxUrlLength   = 4 * 1024
)

// MsgCreatePost defines a CreatePost message
type MsgCreatePost struct {
	UUID         string         `json:"uuid"`
	Owner        sdk.AccAddress `json:"owner"`
	Title        string         `json:"title"`
	PreviewImage string         `json:"preview_image"`
	Text         string         `json:"text"`
	Tags         []string       `json:"tags"`
}

// MsgDeletePost defines a DeletePost message
type MsgDeletePost struct {
	UUID  string         `json:"uuid"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgCreatePost is a constructor function for MsgCreatePost
func NewMsgCreatePost(title string, previewImage string, text string, tags []string, owner sdk.AccAddress) MsgCreatePost {
	return MsgCreatePost{
		UUID:         uuid.Must(uuid.NewV1()).String(),
		Owner:        owner,
		Title:        title,
		PreviewImage: previewImage,
		Text:         text,
		Tags:         tags,
	}
}

// Route should return the name of the module
func (msg MsgCreatePost) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreatePost) Type() string { return "create_post" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreatePost) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if len(msg.Title) > maxTitleLength || len(msg.Title) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "title's length should be less then %d and not zero", maxTitleLength)
	}

	if !IsPreviewImageValid(msg.PreviewImage) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "preview_image is invalid")
	}

	if len(msg.Text) < minPostLength || len(msg.Text) > maxPostLength {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "post's length should be between %d and %d", minPostLength, maxPostLength)
	}

	if len(msg.Tags) > maxTagsCount {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "tags count should be less then %d", maxTagsCount)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreatePost) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreatePost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// NewMsgDeletePost is a constructor function for MsgDeletePost
func NewMsgDeletePost(id uuid.UUID, owner sdk.AccAddress) MsgDeletePost {
	return MsgDeletePost{
		Owner: owner,
		UUID:  id.String(),
	}
}

// Route should return the name of the module
func (msg MsgDeletePost) Route() string { return RouterKey }

// Type should return the action
func (msg MsgDeletePost) Type() string { return "delete_post" }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeletePost) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}

	if _, err := uuid.FromString(msg.UUID); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid uuid")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeletePost) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeletePost) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

func IsPreviewImageValid(str string) bool {
	if len(str) > maxUrlLength {
		return false
	}

	url, err := url.Parse(str)
	if err != nil {
		return false
	}
	return url.Scheme == "http" || url.Scheme == "https"
}
