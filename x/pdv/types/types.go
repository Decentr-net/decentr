package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PDVType int

const (
	PDVTypeCookie PDVType = 1
)

type PDV struct {
	Timestamp time.Time      `json:"timestamp"`
	Address   string         `json:"address"`
	Owner     sdk.AccAddress `json:"owner"`
	Type      PDVType        `json:"type"`
}

// implement fmt.Stringer
func (w PDV) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Timestamp: %s
Owner: %s
Address: %s
Type :%d`, w.Timestamp, w.Owner, w.Address, w.Type))
}
