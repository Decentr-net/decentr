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
func (c *client) SavePDV(ctx context.Context, p schema.PDV) (uint64, error) {
	if !p.Validate() {
		return 0, ErrInvalidRequest
	}

	data, err := json.Marshal(p)
	if err != nil {
		return 0, fmt.Errorf("failed to decode pdv: %w", err)
	}

	data, err = c.sendRequest(ctx, http.MethodPost, "v1/pdv", data)
	if err != nil {
		return 0, fmt.Errorf("failed to make SavePDV request: %w", err)
	}

	resp := SavePDVResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, fmt.Errorf("failed to decode json: %w", err)
	}

	return resp.ID, nil
}

// SavePDV lists pdv for owner.
func (c *client) ListPDV(ctx context.Context, owner string, from uint64, limit uint16) ([]uint64, error) {
	if limit == 0 {
		limit = 100
	}

	data, err := c.sendRequest(ctx, http.MethodGet, fmt.Sprintf("v1/pdv/%s?from=%d&limit=%d", owner, from, limit), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make ListPDV request: %w", err)
	}

	list := make([]uint64, 0)

	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("failed to unmarshal list: %w", err)
	}

	return list, nil
}

// ReceivePDV receives bytes slice from Cerberus by provided address.
func (c *client) ReceivePDV(ctx context.Context, owner string, id uint64) (schema.PDV, error) {
	data, err := c.sendRequest(ctx, http.MethodGet, fmt.Sprintf("v1/pdv/%s/%d", owner, id), nil)
	if err != nil {
		return schema.PDV{}, fmt.Errorf("failed to make ReceivePDV request: %w", err)
	}

	var pdv schema.PDV
	if err := json.Unmarshal(data, &pdv); err != nil {
		return schema.PDV{}, fmt.Errorf("failed to unmarshal pdv: %w", err)
	}

	return pdv, nil
}

// GetPDVMeta returns PDVMeta by provided address.
func (c *client) GetPDVMeta(ctx context.Context, owner string, id uint64) (PDVMeta, error) {
	data, err := c.sendRequest(ctx, http.MethodGet, fmt.Sprintf("v1/pdv/%s/%d/meta", owner, id), nil)
	if err != nil {
		return PDVMeta{}, fmt.Errorf("failed to make GetPDVMeta request: %w", err)
	}

	var m PDVMeta
	if err := json.Unmarshal(data, &m); err != nil {
		return PDVMeta{}, fmt.Errorf("failed to unmarshal meta: %w", err)
	}

	return m, nil
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
