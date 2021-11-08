package config

import (
	"github.com/tendermint/spm/cosmoscmd"
)

const (
	AccountAddressPrefix = "decentr"
	AppName              = "decentr"

	// DefaultBondDenom is the default bond denomination
	DefaultBondDenom = "udec"
)

func SetAddressPrefixes() {
	cosmoscmd.SetPrefixes(AccountAddressPrefix)
}
