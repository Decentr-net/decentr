package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

const DateFormat = "2006-01-02"

func IsValidDate(s string) bool {
	dt, err := time.Parse(DateFormat, s)
	return err == nil && dt.Year() > 1900 && dt.Year() < time.Now().Year()
}

func IsValidGender(str string) bool {
	return str == string(GenderMale) || str == string(GenderFemale)
}

// Public profile data
type Public struct {
	Gender   Gender `json:"gender"`
	Birthday string `json:"birthday"`
}

// Profile represent an account settings storage
type Profile struct {
	// Owner is Profile owner
	Owner sdk.AccAddress `json:"owner"`
	// Public profile data
	Public Public `json:"public"`
	// Encrypted profile data such as emails, phones, nicknames
	Private []byte `json:"private"`
}
