package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PDV struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// implement fmt.Stringer
func (w PDV) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s`, w.Owner, w.Value))
}
