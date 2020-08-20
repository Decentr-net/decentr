package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var ErrNotFound = sdkerrors.Register(ModuleName, 1, "not found")
