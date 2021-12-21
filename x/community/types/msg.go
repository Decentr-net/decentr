package types

import (
	"github.com/gofrs/uuid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewMsgCreatePost is a constructor function for MsgCreatePost
func NewMsgCreatePost(
	title string,
	category Category,
	previewImage string,
	text string,
	owner sdk.AccAddress,
) MsgCreatePost {
	return MsgCreatePost{
		Post: Post{
			Uuid:         uuid.Must(uuid.NewV1()).String(),
			Owner:        owner.String(),
			Title:        title,
			Category:     category,
			PreviewImage: previewImage,
			Text:         text,
		},
	}
}

// Route should return the name of the module
func (m MsgCreatePost) Route() string { return RouterKey }

// Type should return the action
func (m MsgCreatePost) Type() string { return "create_post" }

// ValidateBasic runs stateless checks on the message
func (m MsgCreatePost) ValidateBasic() error {
	if err := m.Post.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgCreatePost) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgCreatePost) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Post.Owner)
	return []sdk.AccAddress{addr}
}

// NewMsgDeletePost is a constructor function for MsgDeletePost
func NewMsgDeletePost(owner, postOwner sdk.AccAddress, postUUID uuid.UUID) MsgDeletePost {
	return MsgDeletePost{
		Owner:     owner.String(),
		PostUuid:  postUUID.String(),
		PostOwner: postOwner.String(),
	}
}

// Route should return the name of the module
func (m MsgDeletePost) Route() string { return RouterKey }

// Type should return the action
func (m MsgDeletePost) Type() string { return "delete_post" }

// ValidateBasic runs stateless checks on the message
func (m MsgDeletePost) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if _, err := uuid.FromString(m.PostUuid); err != nil {
		return sdkerrors.ErrInvalidRequest.Wrapf("invalid post_uuid: %s", err)
	}

	postOwner, err := sdk.AccAddressFromBech32(m.PostOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid post_owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(postOwner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid post_owner address: %s", err)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgDeletePost) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgDeletePost) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

// NewMsgSetLike is a constructor function for MsgSetLike
func NewMsgSetLike(postOwner sdk.AccAddress, postUUID uuid.UUID, owner sdk.AccAddress, weight LikeWeight) MsgSetLike {
	return MsgSetLike{
		Like: Like{
			PostOwner: postOwner.String(),
			PostUuid:  postUUID.String(),
			Owner:     owner.String(),
			Weight:    weight,
		},
	}
}

// Route should return the name of the module
func (m MsgSetLike) Route() string { return RouterKey }

// Type should return the action
func (m MsgSetLike) Type() string { return "set_like" }

// ValidateBasic runs stateless checks on the message
func (m MsgSetLike) ValidateBasic() error {
	if err := m.Like.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgSetLike) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgSetLike) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Like.Owner)
	return []sdk.AccAddress{addr}
}

// NewMsgFollow is a constructor function for MsgFollow
func NewMsgFollow(owner, whom sdk.AccAddress) MsgFollow {
	return MsgFollow{
		Owner: owner.String(),
		Whom:  whom.String(),
	}
}

// Route should return the name of the module
func (m MsgFollow) Route() string { return RouterKey }

// Type should return the action
func (m MsgFollow) Type() string { return "follow" }

// ValidateBasic runs stateless checks on the message
func (m MsgFollow) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	whom, err := sdk.AccAddressFromBech32(m.Whom)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid whom address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(whom); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid whom address: %s", err)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgFollow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgFollow) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}

// NewMsgUnfollow is a constructor function for MsgUnfollow
func NewMsgUnfollow(owner, whom sdk.AccAddress) MsgUnfollow {
	return MsgUnfollow{
		Owner: owner.String(),
		Whom:  whom.String(),
	}
}

// Route should return the name of the module
func (m MsgUnfollow) Route() string { return RouterKey }

// Type should return the action
func (m MsgUnfollow) Type() string { return "unfollow" }

// ValidateBasic runs stateless checks on the message
func (m MsgUnfollow) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	whom, err := sdk.AccAddressFromBech32(m.Whom)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(whom); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid whom address: %s", err)
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgUnfollow) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

// GetSigners defines whose signature is required
func (m MsgUnfollow) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Owner)
	return []sdk.AccAddress{addr}
}
