// Package schema provides schemas and validation functions for it.
package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// PDVVersion represents version.
// swagger:enum PDVVersion
type PDVVersion string

const (
	// PDVV1 ...
	PDVV1 PDVVersion = "v1"
)

const (
	// PDVDataSizeLimit is limit to PDVData's size.
	PDVDataSizeLimit = 8 * 1024
)

// nolint: gochecknoglobals
var (
	pdvObjectSchemes = map[PDVVersion]reflect.Type{
		PDVV1: reflect.TypeOf(PDVObjectV1{}),
	}
)

// PDV is main Data object.
type PDV struct {
	Version PDVVersion `json:"version"`

	PDV []PDVObject `json:"pdv"`
}

// PDVObject is interface for all versions objects.
type PDVObject interface {
	Validate

	GetData() []PDVData
}

// PDVObjectMetaV1 is PDVObjectV1 meta Data.
type PDVObjectMetaV1 struct {
	// Domain of website where object was taken
	Host string `json:"domain"`
	// Path of website's url where object was taken
	Path string `json:"path"`
}

// PDVObjectV1 is PDVObject implementation with v1 version.
type PDVObjectV1 struct {
	PDVObjectMetaV1

	Data []PDVData `json:"data"`
}

// GetData returns slice of PDVData.
func (o PDVObjectV1) GetData() []PDVData {
	return o.Data
}

// UnmarshalJSON ...
func (p *PDV) UnmarshalJSON(b []byte) error {
	var i struct {
		Version PDVVersion `json:"version"`

		PDV []json.RawMessage `json:"pdv"`
	}

	if err := json.Unmarshal(b, &i); err != nil {
		return fmt.Errorf("failed to unmarshal PDV meta: %w", err)
	}

	t, ok := pdvObjectSchemes[i.Version]
	if !ok {
		return errors.New("unknown version of object")
	}

	p.Version = i.Version

	for _, v := range i.PDV {
		o := reflect.New(t).Interface()

		if err := json.Unmarshal(v, o); err != nil {
			return err
		}

		p.PDV = append(p.PDV, o.(PDVObject))
	}

	return nil
}

// UnmarshalJSON ...
func (o *PDVObjectV1) UnmarshalJSON(b []byte) error {
	var i struct {
		PDVObjectMetaV1

		Data []json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*o = PDVObjectV1{
		PDVObjectMetaV1: i.PDVObjectMetaV1,
		Data:            make([]PDVData, len(i.Data)),
	}

	for i, v := range i.Data {
		if len(v) > PDVDataSizeLimit {
			return errors.New("pdv Data is too big")
		}

		type T struct {
			Type PDVType `json:"type"`
		}

		var m T
		if err := json.Unmarshal(v, &m); err != nil {
			return fmt.Errorf("failed to unmarshal PDV Data meta: %w", err)
		}

		t, ok := pdvDataSchemes[m.Type]
		if !ok {
			return fmt.Errorf("unknown pdv Data: %s", m.Type)
		}

		d := reflect.New(t).Interface().(PDVData) // nolint:errcheck

		if err := json.Unmarshal(v, d); err != nil {
			return fmt.Errorf("failed to unmarshal Data: %w", err)
		}

		o.Data[i] = d
	}

	return nil
}
