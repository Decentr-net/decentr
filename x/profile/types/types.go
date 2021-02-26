package types

import (
	"net/url"
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

func IsValidAvatar(str string) bool {
	if len(str) > 4*1024 {
		return false
	}

	url, err := url.Parse(str)
	if err != nil {
		return false
	}
	return url.Scheme == "http" || url.Scheme == "https"
}

// Public profile data
type Public struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Avatar       string `json:"avatar"`
	Bio          string `json:"bio"`
	Gender       Gender `json:"gender"`
	Birthday     string `json:"birthday"`
	RegisteredAt int64  `json:"registeredAt"`
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
