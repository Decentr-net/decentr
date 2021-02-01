package schema

import (
	"encoding/json"
	"errors"
	"reflect"
)

// PDVType represents Data type.
// swagger:enum PDVType
type PDVType string

const (
	// PDVCookieType ...
	PDVCookieType PDVType = "cookie"
	// PDVLoginCookieType ...
	PDVLoginCookieType PDVType = "login_cookie"
)

var pdvDataSchemes = map[PDVType]reflect.Type{ // nolint:gochecknoglobals
	PDVCookieType:      reflect.TypeOf(PDVDataCookie{}),
	PDVLoginCookieType: reflect.TypeOf(PDVDataLoginCookie{}),
}

// PDVData is interface for all Data types.
type PDVData interface {
	Validate

	Type() PDVType
}

// PDVDataCookie is PDVData implementation for Cookies(according to https://developer.chrome.com/extensions/cookies).
type PDVDataCookie struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	Domain         string `json:"domain"`
	Path           string `json:"path"`
	SameSite       string `json:"same_site"`
	HostOnly       bool   `json:"host_only"`
	Secure         bool   `json:"secure"`
	ExpirationDate uint64 `json:"expiration_date,omitempty"`
}

// Type ...
func (PDVDataCookie) Type() PDVType {
	return PDVCookieType
}

// Validate ...
func (d *PDVDataCookie) Validate() bool {
	if d.Name == "" || d.Value == "" {
		return false
	}

	return true
}

// MarshalJSON ...
func (d PDVDataCookie) MarshalJSON() ([]byte, error) { // nolint:gocritic
	type T PDVDataCookie
	v := struct {
		Type PDVType `json:"type"`
		T
	}{
		Type: d.Type(),
		T:    T(d),
	}

	return json.Marshal(v)
}

// PDVDataLoginCookie is the same as PDVDataCookie but with different type.
type PDVDataLoginCookie struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	Domain         string `json:"domain"`
	Path           string `json:"path"`
	SameSite       string `json:"same_site"`
	HostOnly       bool   `json:"host_only"`
	Secure         bool   `json:"secure"`
	ExpirationDate uint64 `json:"expiration_date,omitempty"`
}

// Type ...
func (PDVDataLoginCookie) Type() PDVType {
	return PDVLoginCookieType
}

// MarshalJSON ...
func (d PDVDataLoginCookie) MarshalJSON() ([]byte, error) { // nolint:gocritic
	type T PDVDataLoginCookie
	v := struct {
		Type PDVType `json:"type"`
		T
	}{
		Type: d.Type(),
		T:    T(d),
	}

	return json.Marshal(v)
}

// Validate ...
func (d *PDVDataLoginCookie) Validate() bool {
	return (*PDVDataCookie)(d).Validate()
}

// UnmarshalText ...
func (t *PDVType) UnmarshalText(b []byte) error {
	s := PDVType(b)
	switch s {
	case PDVCookieType, PDVLoginCookieType:
	default:
		return errors.New("unknown PDVType")
	}
	*t = s
	return nil
}
