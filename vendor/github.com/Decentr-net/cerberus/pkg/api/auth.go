package api

import (
	"bytes"
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

	if len(k) != 33 {
		return nil, nil, ErrInvalidPublicKey
	}

	pk, err := amino.PubKeyFromBytes(GetAminoSecp256k1PubKey(k))
	if err != nil {
		return nil, nil, ErrInvalidPublicKey
	}

	return pk, s, nil
}

// Sign signs http request.
func Sign(r *http.Request, pk crypto.PrivKey) error {
	d, err := GetMessageToSign(r)
	if err != nil {
		return fmt.Errorf("failed to get digest: %w", err)
	}

	s, err := pk.Sign(d)
	if err != nil {
		return fmt.Errorf("failed to sign digest: %w", err)
	}

	r.Header.Set(PublicKeyHeader, hex.EncodeToString(pk.PubKey().Bytes()[5:])) // truncate amino codec prefix
	r.Header.Set(SignatureHeader, hex.EncodeToString(s))

	return nil
}

// Verify verifies request's signature with public key.
func Verify(r *http.Request) error {
	k, s, err := getSignature(r)
	if err != nil {
		return err
	}

	d, err := GetMessageToSign(r)
	if err != nil {
		return err
	}
	if !k.VerifyBytes(d, s) {
		return ErrNotVerified
	}

	return nil
}

// GetMessageToSign returns message to sign.
func GetMessageToSign(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	return append(body, []byte(r.URL.Path)...), nil
}

// GetAminoSecp256k1PubKey adds amino secp256k1 pubkey prefix to pubkey(including length-prefix).
func GetAminoSecp256k1PubKey(k []byte) []byte {
	return append([]byte{235, 90, 233, 135, 33}, k...)
}
