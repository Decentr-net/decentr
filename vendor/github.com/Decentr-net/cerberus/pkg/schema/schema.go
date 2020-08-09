// Package schema provides schemas and validation functions for it.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const (
	// PDVv1 ...
	PDVv1 PDVVersion = "v1"
)
const (
	// PDVCookieType ...
	PDVCookieType PDVType = "cookie"
)
const (
	// DataPerPDVLimit is limit how much data we allow per PDV.
	DataPerPDVLimit = 20
	// PDVDataSizeLimit is limit to PDVData's size.
	PDVDataSizeLimit = 8 * 1024
)

// PDVVersion represents version.
type PDVVersion string

// PDVType represents data type.
type PDVType string

// nolint: gochecknoglobals
var (
	pdvObjectSchemes = map[PDVVersion]reflect.Type{
		PDVv1: reflect.TypeOf(PDVObjectV1{}),
	}

	pdvDataSchemes = map[PDVType]map[PDVVersion]reflect.Type{
		PDVCookieType: {
			PDVv1: reflect.TypeOf(PDVDataCookieV1{}),
		},
	}
)

// PDV is main data object.
type PDV struct {
	Version PDVVersion `json:"version"`

	PDV PDVObject `json:"pdv"`
}

// PDVObject is interface for all versions objects.
type PDVObject interface {
	Validate

	Version() PDVVersion
}

// PDVObjectMetaV1 is PDVObjectV1 meta data.
type PDVObjectMetaV1 struct {
	// Website information
	Host string `json:"domain"`
	Path string `json:"path"`
}

// PDVObjectV1 is PDVObject implementation with v1 version.
type PDVObjectV1 struct {
	PDVObjectMetaV1

	Data []PDVData `json:"data"`
}

// Version ...
func (o *PDVObjectV1) Version() PDVVersion {
	return PDVv1
}

// PDVDataMeta contains common information about data.
type PDVDataMeta struct {
	PDVVersion PDVVersion `json:"version"`
	PDVType    PDVType    `json:"type"`
}

// PDVData is interface for all data types.
type PDVData interface {
	Validate

	Version() PDVVersion
	Type() PDVType
}

// PDVDataCookieV1 is PDVData implementation for Cookies(according to https://developer.chrome.com/extensions/cookies) with version v1.
type PDVDataCookieV1 struct {
	Name           string `json:"name"`
	Value          string `json:"value"`
	Domain         string `json:"domain"`
	Path           string `json:"path"`
	SameSite       string `json:"same_site"`
	HostOnly       bool   `json:"host_only"`
	Secure         bool   `json:"secure"`
	ExpirationDate uint64 `json:"expiration_date,omitempty"`
}

// UnmarshalJSON ...
func (p *PDV) UnmarshalJSON(b []byte) error {
	var i struct {
		Version PDVVersion `json:"version"`

		PDV json.RawMessage `json:"pdv"`
	}

	if err := json.Unmarshal(b, &i); err != nil {
		return fmt.Errorf("failed to unmarshal PDV meta: %w", err)
	}

	t, ok := pdvObjectSchemes[i.Version]
	if !ok {
		return errors.New("unknown version of object")
	}

	v := reflect.New(t).Interface()
	if err := json.Unmarshal(i.PDV, v); err != nil {
		return err
	}

	p.Version = i.Version
	p.PDV = v.(PDVObject) // nolint

	return nil
}

// UnmarshalJSON ...
func (o *PDVObjectV1) UnmarshalJSON(b []byte) error {
	var i struct {
		PDVObjectMetaV1

		PDVData []json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	if len(i.PDVData) > DataPerPDVLimit {
		return errors.New("too much data in PDV")
	}

	*o = PDVObjectV1{
		PDVObjectMetaV1: i.PDVObjectMetaV1,
		Data:            make([]PDVData, len(i.PDVData)),
	}

	for i, v := range i.PDVData {
		if len(v) > PDVDataSizeLimit {
			return errors.New("pdv data is too big")
		}

		var m PDVDataMeta
		if err := json.Unmarshal(v, &m); err != nil {
			return fmt.Errorf("failed to unmarshal PDV data meta: %w", err)
		}

		t, ok := pdvDataSchemes[m.PDVType][m.PDVVersion]
		if !ok {
			return fmt.Errorf("unknown pdv data: %s %s", m.PDVType, m.PDVVersion)
		}

		d := reflect.New(t).Interface().(PDVData) // nolint:errcheck

		if err := json.Unmarshal(v, d); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		o.Data[i] = d
	}

	return nil
}

// Version ...
func (PDVDataCookieV1) Version() PDVVersion {
	return PDVv1
}

// Type ...
func (PDVDataCookieV1) Type() PDVType {
	return PDVCookieType
}

// MarshalJSON ...
func (p PDV) MarshalJSON() ([]byte, error) {
	p.Version = p.PDV.Version()
	type t PDV
	return json.Marshal(t(p))
}

// MarshalJSON ...
func (d PDVDataCookieV1) MarshalJSON() ([]byte, error) { // nolint:gocritic
	type T PDVDataCookieV1
	v := struct {
		PDVDataMeta
		T
	}{
		PDVDataMeta: PDVDataMeta{
			PDVVersion: d.Version(),
			PDVType:    d.Type(),
		},
		T: T(d),
	}

	return json.Marshal(v)
}
