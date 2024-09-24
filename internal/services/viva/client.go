package viva

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"safe-size-pay/internal/constants"
)

// VivaClient interface defines the methods that our client should implement
type VivaClient interface {
	GetToken() (*Token, error)
	CreateOrder(request *OrderRequest) (*OrderResponse, error)
}

// client struct holds the configuration and HTTP client
type client struct {
	accountsUrl string
	apiUrl      string
	username    string
	password    string
	httpClient  *http.Client
}

// NewClient creates a new instance of client
func NewClient() VivaClient {
	return &client{
		accountsUrl: os.Getenv(constants.VivaBaseAccountsUrl),
		apiUrl:      os.Getenv(constants.VivaBaseApiUrl),
		username:    os.Getenv(constants.VivaUsername),
		password:    os.Getenv(constants.VivaPassword),
		httpClient:  &http.Client{},
	}
}

// GetToken implements the GetToken method of the VivaClient interface
func (vc *client) GetToken() (*Token, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/connect/token", vc.accountsUrl), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", constants.FormUrlEncoded)

	auth := vc.username + ":" + vc.password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, fmt.Errorf("error parsing token response: %v", err)
	}
	return &token, nil
}

// CreateOrder implements the CreateOrder method of the Client interface
func (vc *client) CreateOrder(request *OrderRequest) (*OrderResponse, error) {
	token, err := vc.GetToken()
	if err != nil {
		return nil, fmt.Errorf("error getting token: %v", err)
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/checkout/v2/orders", vc.apiUrl), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", constants.ApplicationJson)
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var orderResponse OrderResponse
	err = json.Unmarshal(body, &orderResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing order response: %v", err)
	}

	return &orderResponse, nil
}
