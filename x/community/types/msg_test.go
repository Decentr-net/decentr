package types

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"

	. "github.com/Decentr-net/decentr/testutil"
)

func TestMsgCreatePost_ValidateBasic(t *testing.T) {
	// check if post validation called
	require.Error(t, MsgCreatePost{}.ValidateBasic())
}

func TestMsgResetAccount_ValidateBasic(t *testing.T) {
	valid := NewMsgDeletePost(NewAccAddress(), NewAccAddress(), uuid.Must(uuid.NewV1()))

	alter := func(f func(m *MsgDeletePost)) MsgDeletePost {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgDeletePost) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDeletePost) {
		m.PostOwner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDeletePost) {
		m.PostUuid = ""
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgDeletePost) {
		m.PostUuid = "123"
	}).ValidateBasic())
}

func TestMsgSetLike_ValidateBasic(t *testing.T) {
	// check if like validation called
	require.Error(t, MsgSetLike{}.ValidateBasic())
}

func TestMsgFollow_ValidateBasic(t *testing.T) {
	valid := NewMsgFollow(NewAccAddress(), NewAccAddress())

	alter := func(f func(m *MsgFollow)) MsgFollow {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgFollow) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgFollow) {
		m.Whom = nil
	}).ValidateBasic())
}

func TestMsgUnfollow_ValidateBasic(t *testing.T) {
	valid := NewMsgUnfollow(NewAccAddress(), NewAccAddress())

	alter := func(f func(m *MsgUnfollow)) MsgUnfollow {
		cp := valid
		f(&cp)
		return cp
	}

	require.NoError(t, valid.ValidateBasic())
	require.Error(t, alter(func(m *MsgUnfollow) {
		m.Owner = nil
	}).ValidateBasic())
	require.Error(t, alter(func(m *MsgUnfollow) {
		m.Whom = nil
	}).ValidateBasic())
}
