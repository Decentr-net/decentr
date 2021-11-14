package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/spm/cosmoscmd"
)

const (
	AccountAddressPrefix = "decentr"
	AppName              = "decentr"

	// DefaultBondDenom is the default bond denomination
	DefaultBondDenom = "udec"
)

var (
	InitialTokenBalance = sdk.NewDec(1)
)

var isConfigSealed bool

func SetAddressPrefixes() {
	if !isConfigSealed {
		cosmoscmd.SetPrefixes(AccountAddressPrefix)
	}
	isConfigSealed = true
}
