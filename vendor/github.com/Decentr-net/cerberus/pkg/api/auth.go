package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tendermint/tendermint/crypto"
	amino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

// PublicKeyHeader is name for public key http header.
const PublicKeyHeader = "Public-Key"

// SignatureHeader is name for signature http header.
const SignatureHeader = "Signature"

func getSignature(r *http.Request) (crypto.PubKey, []byte, error) {
	s, err := hex.DecodeString(r.Header.Get(SignatureHeader))

	if err != nil {
		return nil, nil, ErrInvalidSignature
	}

	k, err := hex.DecodeString(r.Header.Get(PublicKeyHeader))
	if err != nil {
		return nil, nil, ErrInvalidPublicKey
	}

	pk, err := amino.PubKeyFromBytes(k)

	if err != nil {
		return nil, nil, ErrInvalidPublicKey
	}

	return pk, s, nil
}

// Sign signs http request.
func Sign(r *http.Request, pk crypto.PrivKey) error {
	d, err := Digest(r)
	if err != nil {
		return fmt.Errorf("failed to get digest: %w", err)
	}

	s, err := pk.Sign(d)
	if err != nil {
		return fmt.Errorf("failed to sign digest: %w", err)
	}

	r.Header.Set(PublicKeyHeader, hex.EncodeToString(pk.PubKey().Bytes()))
	r.Header.Set(SignatureHeader, hex.EncodeToString(s))

	return nil
}

// Verify verifies request's signature with public key.
func Verify(r *http.Request) ([]byte, error) {
	k, s, err := getSignature(r)
	if err != nil {
		return nil, err
	}

	d, err := Digest(r)
	if err != nil {
		return nil, err
	}
	if !k.VerifyBytes(d, s) {
		return nil, ErrNotVerified
	}

	return d, nil
}

// Digest returns sha256 of request body.
func Digest(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	body = append(body, []byte(r.URL.Path)...)

	d := sha256.Sum256(body)
	return d[:], nil
}
