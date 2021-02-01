// Package api provides API and client to Cerberus.
package api

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/Decentr-net/cerberus/pkg/schema"
)

//go:generate mockgen -destination=./api_mock.go -package=api -source=api.go

// nolint: gochecknoglobals
var addressRegExp = regexp.MustCompile(`[0-9a-fA-F]{40}-[0-9a-fA-F]{64}`) // public_key_hex/data_sha256_digest_hex

// ErrInvalidRequest is returned when request is invalid.
var ErrInvalidRequest = errors.New("invalid request")

// ErrInvalidPublicKey is returned when public key is invalid.
var ErrInvalidPublicKey = fmt.Errorf("%w: public key is invalid", ErrInvalidRequest)

// ErrInvalidSignature is returned when signature is invalid.
var ErrInvalidSignature = fmt.Errorf("%w: signature is invalid", ErrInvalidRequest)

// ErrNotFound is returned when object is not found.
var ErrNotFound = errors.New("not found")

// ErrNotVerified is returned when signature is wrong.
var ErrNotVerified = errors.New("failed to verify message")

// Cerberus provides user-friendly API methods.
type Cerberus interface {
	SavePDV(ctx context.Context, p schema.PDV) (string, error)
	ReceivePDV(ctx context.Context, address string) (schema.PDV, error)
	GetPDVMeta(ctx context.Context, address string) (PDVMeta, error)
}

// Error ...
// swagger:model Error
type Error struct {
	Error string `json:"error"`
}

// SavePDVResponse ...
// swagger:model SavePDVResponse
type SavePDVResponse struct {
	Address string `json:"address"`
}

// PDVMeta contains info about PDV.
type PDVMeta struct {
	// ObjectTypes represents how much certain pdv data pdv contains.
	ObjectTypes map[schema.PDVType]uint16 `json:"object_types"`
	Reward      uint64                    `json:"reward"`
}

// IsAddressValid check is address is matching with regexp.
func IsAddressValid(s string) bool {
	return addressRegExp.MatchString(s)
}
