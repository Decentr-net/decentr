package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/Decentr-net/cerberus/pkg/schema"
)

type client struct {
	host string

	pk secp256k1.PrivKeySecp256k1

	c *http.Client
}

// NewClient returns client with http.DefaultClient.
func NewClient(host string, pk secp256k1.PrivKeySecp256k1) Cerberus {
	return NewClientWithHTTPClient(host, pk, &http.Client{})
}

// NewClientWithHTTPClient returns client with provided http.Client.
func NewClientWithHTTPClient(host string, pk secp256k1.PrivKeySecp256k1, c *http.Client) Cerberus {
	return &client{
		host: host,
		pk:   pk,
		c:    c,
	}
}

// SavePDV sends bytes slice to Cerberus.
// SavePDV can return ErrInvalidRequest besides general api package's errors.
func (c *client) SavePDV(ctx context.Context, p *schema.PDV) (string, error) {
	// validate data

	if !p.PDV.Validate() {
		return "", ErrInvalidRequest
	}

	data, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("failed to decode pdv: %w", err)
	}

	data, err = c.sendRequest(ctx, http.MethodPost, "v1/pdv", data)
	if err != nil {
		return "", fmt.Errorf("failed to make SavePDV request: %w", err)
	}

	resp := SavePDVResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("failed to decode json: %w", err)
	}

	return resp.Address, nil
}

// ReceivePDV receives bytes slice from Cerberus by provided address.
// ReceivePDV can return ErrInvalidRequest and ErrNotFound besides general api package's errors.
func (c *client) ReceivePDV(ctx context.Context, address string) (json.RawMessage, error) {
	if !IsAddressValid(address) {
		return nil, ErrInvalidRequest
	}

	data, err := c.sendRequest(ctx, http.MethodGet, fmt.Sprintf("v1/pdv/%s", address), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make ReceivePDV request: %w", err)
	}

	return data, nil
}

// DoesPDVExist returns is data exists in Cerberus by provided address.
// DoesPDVExist can return ErrInvalidRequest and ErrNotFound besides general api package's errors.
func (c *client) DoesPDVExist(ctx context.Context, address string) (bool, error) {
	if !IsAddressValid(address) {
		return false, ErrInvalidRequest
	}

	_, err := c.sendRequest(ctx, http.MethodHead, fmt.Sprintf("v1/pdv/%s", address), nil)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to make DoesPDVExist request: %w", err)
	}

	return true, nil
}

// sendRequest is utility method which signs request, if it's needed, and send POST request to Cerberus.
// Also converts http.StatusCode to package's errors.
func (c *client) sendRequest(ctx context.Context, method, endpoint string, body []byte) ([]byte, error) {
	r, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s/%s", c.host, endpoint), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err := Sign(r, c.pk); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	resp, err := c.c.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, ErrNotFound
		case http.StatusBadRequest:
			return nil, ErrInvalidRequest
		default:
			var e Error
			if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
				return nil, errors.Errorf("request failed with status %d", resp.StatusCode)
			}
			return nil, errors.Errorf("request failed: %s", e.Error)
		}
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}
