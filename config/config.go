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

func SetAddressPrefixes() {
	cosmoscmd.SetPrefixes(AccountAddressPrefix)
}
