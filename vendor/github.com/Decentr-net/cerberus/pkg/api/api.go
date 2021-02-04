// Package api provides API and client to Cerberus.
package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/Decentr-net/cerberus/pkg/schema"
)

//go:generate mockgen -destination=./api_mock.go -package=api -source=api.go

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
	SavePDV(ctx context.Context, p schema.PDV) (uint64, error)
	ListPDV(ctx context.Context, owner string, from uint64, limit uint16) ([]uint64, error)
	ReceivePDV(ctx context.Context, owner string, id uint64) (schema.PDV, error)
	GetPDVMeta(ctx context.Context, owner string, id uint64) (PDVMeta, error)
}

// Error ...
// swagger:model Error
type Error struct {
	Error string `json:"error"`
}

// SavePDVResponse ...
// swagger:model SavePDVResponse
type SavePDVResponse struct {
	ID uint64 `json:"id"`
}

// PDVMeta contains info about PDV.
type PDVMeta struct {
	// ObjectTypes represents how much certain pdv data pdv contains.
	ObjectTypes map[schema.PDVType]uint16 `json:"object_types"`
	Reward      uint64                    `json:"reward"`
}
