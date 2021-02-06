package types

import (
	"net/url"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gofrs/uuid"
)

const (
	maxTitleLength = 150
	maxPostLength  = 64 * 1000
	minPostLength  = 15
	maxURLLength   = 4 * 1024
)

// MsgCreatePost defines a CreatePost message
type MsgCreatePost struct {
	UUID         string         `json:"uuid"`
	Owner        sdk.AccAddress `json:"owner"`
	Title        string         `json:"title"`
	Category     Category       `json:"category"`
	PreviewImage string         `json:"previewImage"`
	Text         string         `json:"text"`
}

// MsgDeletePost defines a DeletePost message
type MsgDeletePost struct {
	PostUUID  string         `json:"postUUID"`
	PostOwner sdk.AccAddress `json:"postOwner"`
	Owner     sdk.AccAddress `json:"owner"`
}

// MsgSetLike defines a SetLike message
type MsgSetLike struct {
	PostOwner sdk.AccAddress `json:"postOwner"`
	PostUUID  string         `json:"postUuid"`
	Owner     sdk.AccAddress `json:"owner"`
	Weight    LikeWeight     `json:"weight"`
}

// MsgFollow defines a MsgFollow message
type MsgFollow struct {
	Owner sdk.AccAddress `json:"owner"`
	Whom  sdk.AccAddress `json:"whom"`
}

// MsgUnfollow defines a MsgUnfollow message
type MsgUnfollow struct {
	Owner sdk.AccAddress `json:"owner"`
	Whom  sdk.AccAddress `json:"whom"`
}

// NewMsgCreatePost is a constructor function for MsgCreatePost
func NewMsgCreatePost(title string, category Category, previewImage string, text string, owner sdk.AccAddress) MsgCreatePost {
	return MsgCreatePost{
		UUID:         uuid.Must(uuid.NewV1()).String(),
		Owner:        owner,
		Title:        title,
		Category:     category,
		PreviewImage: previewImage,
		Text:         text,
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "title should be shorter then %d and not empty", maxTitleLength)
	}

	if msg.Category == UndefinedCategory || msg.Category > SportsCategory {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid category")
	}

	if !IsPreviewImageValid(msg.PreviewImage) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "previewImage is invalid")
	}

	if utf8.RuneCountInString(msg.Text) < minPostLength || len(msg.Text) > maxPostLength {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "post's length should be between %d symbols and %d bytes", minPostLength, maxPostLength)
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
func NewMsgDeletePost(owner sdk.AccAddress, postUUID uuid.UUID, postOwner sdk.AccAddress) MsgDeletePost {
	return MsgDeletePost{
		Owner:     owner,
		PostUUID:  postUUID.String(),
		PostOwner: postOwner,
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

	if msg.PostOwner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.PostOwner.String())
	}

	if _, err := uuid.FromString(msg.PostUUID); err != nil {
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
	if len(str) > maxURLLength {
		return false
	}

	if str == "" {
		return true
	}

	url, err := url.Parse(str)
	if err != nil {
		return false
	}
	return url.Scheme == "http" || url.Scheme == "https"
}

// NewMsgSetLike is a constructor function for MsgSetLike
func NewMsgSetLike(postOwner sdk.AccAddress, postUUID uuid.UUID, owner sdk.AccAddress, weight LikeWeight) MsgSetLike {
	return MsgSetLike{
		PostOwner: postOwner,
		PostUUID:  postUUID.String(),
		Owner:     owner,
		Weight:    weight,
	}
}

// Route should return the name of the module
func (msg MsgSetLike) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetLike) Type() string { return "set_like" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetLike) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner is empty")
	}

	if _, err := uuid.FromString(msg.PostUUID); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid post uuid")
	}

	if msg.Weight > LikeWeightUp || msg.Weight < LikeWeightDown {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid weight")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetLike) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetLike) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// NewMsgFollow is a constructor function for MsgFollow
func NewMsgFollow(owner, whom sdk.AccAddress) MsgFollow {
	return MsgFollow{
		Owner: owner,
		Whom:  whom,
	}
}

// Route should return the name of the module
func (msg MsgFollow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgFollow) Type() string { return "follow" }

// ValidateBasic runs stateless checks on the message
func (msg MsgFollow) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner is empty")
	}

	if msg.Whom.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "whom is empty")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgFollow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgFollow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// NewMsgUnfollow is a constructor function for MsgUnfollow
func NewMsgUnfollow(owner, whom sdk.AccAddress) MsgUnfollow {
	return MsgUnfollow{
		Owner: owner,
		Whom:  whom,
	}
}

// Route should return the name of the module
func (msg MsgUnfollow) Route() string { return RouterKey }

// Type should return the action
func (msg MsgUnfollow) Type() string { return "unfollow" }

// ValidateBasic runs stateless checks on the message
func (msg MsgUnfollow) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner is empty")
	}

	if msg.Whom.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "whom is empty")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgUnfollow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgUnfollow) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
