package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PDV struct {
	Address string         `json:"address"`
	Owner   sdk.AccAddress `json:"owner"`
}

// implement fmt.Stringer
func (w PDV) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Address: %s`, w.Owner, w.Address))
}
