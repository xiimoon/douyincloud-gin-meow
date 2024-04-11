package oceanApi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type APIClient struct {
	BaseURL       string
	HTTPClient    *http.Client
	AppID         string
	Secret        string
	AccessToken   string
	mu            sync.Mutex
	refreshTicker *time.Ticker
}

type APIResponse struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}

func NewAPIClient(baseURL, appID, secret string) *APIClient {
	client := &APIClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
		AppID:      appID,
		Secret:     secret,
	}
	client.refreshToken() // 初始化时获取token
	return client
}

func (client *APIClient) refreshToken() {
	client.mu.Lock()
	defer client.mu.Unlock()

	// Initial backoff interval is 1 second.
	backoffInterval := 1 * time.Second

	for {
		requestData := map[string]string{
			"appid":      client.AppID,
			"secret":     client.Secret,
			"grant_type": "client_credential",
		}

		jsonData, err := json.Marshal(requestData)
		if err != nil {
			log.Printf("Error marshalling token request data: %v", err)
			return
		}

		req, err := http.NewRequest("POST", client.BaseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error creating token request: %v", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.HTTPClient.Do(req)
		if err != nil {
			log.Printf("Error sending token request: %v", err)
			time.Sleep(backoffInterval)
			backoffInterval *= 2 // Double the backoff interval for the next retry.
			continue
		}
		defer resp.Body.Close()

		var apiResponse APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
			log.Printf("Error decoding token response: %v", err)
			return
		}

		if apiResponse.ErrNo != 0 {
			log.Printf("API error: %s (%d)", apiResponse.ErrTips, apiResponse.ErrNo)
			time.Sleep(backoffInterval)
			backoffInterval *= 2 // Double the backoff interval for the next retry.
			continue
		}

		client.AccessToken = apiResponse.Data.AccessToken

		expiresDuration := time.Duration(apiResponse.Data.ExpiresIn-300) * time.Second
		if client.refreshTicker != nil {
			client.refreshTicker.Stop()
		}
		client.refreshTicker = time.NewTicker(expiresDuration)
		go func() {
			for range client.refreshTicker.C {
				client.refreshToken()
			}
		}()

		break // Token refreshed successfully, exit the loop.
	}
}

// Here you would add your Get, Post, or other HTTP method functions
// that utilize the AccessToken for authentication with the API.
// Get sends an HTTP GET request and parses the JSON response.
func (client *APIClient) Get(endpoint string, response interface{}) error {
	req, err := http.NewRequest("GET", client.BaseURL+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+client.AccessToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	return nil
}

// Post sends an HTTP POST request with JSON data and parses the JSON response.
func (client *APIClient) Post(endpoint string, requestData interface{}, response interface{}) error {
	data, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", client.BaseURL+endpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.AccessToken)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	return nil
}
