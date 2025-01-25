package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AssetPortal/assets-api/pkg/model"
)

type PolkadotClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPolkadotClient(baseURL string, httpClient *http.Client) *PolkadotClient {
	return &PolkadotClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *PolkadotClient) VerifySignature(ctx context.Context, message, address, signature string) (*model.Auth, error) {
	payload := map[string]string{
		"message":   message,
		"address":   address,
		"signature": signature,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/verify", c.baseURL),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	var auth model.Auth
	err = json.NewDecoder(resp.Body).Decode(&auth)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &auth, nil
}
