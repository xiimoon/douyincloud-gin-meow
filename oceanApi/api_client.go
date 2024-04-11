package oceanApi

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type APIClient struct {
	BaseURL              string
	HTTPClient           *http.Client
	AppID                string
	Secret               string
	PublicKeyFingerprint string
	AccessToken          string
	mu                   sync.Mutex
	refreshTicker        *time.Ticker
}

type APIResponse struct {
	ErrNo   int    `json:"err_no"`
	ErrTips string `json:"err_tips"`
	Data    struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}

func NewAPIClient(baseURL, appID, secret, publicKeyFingerprint string) *APIClient {
	client := &APIClient{
		BaseURL:              baseURL,
		AppID:                appID,
		Secret:               secret,
		PublicKeyFingerprint: publicKeyFingerprint,
	}
	client.configureHTTPClient()
	client.refreshToken()
	return client
}

func (client *APIClient) configureHTTPClient() {
	client.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 注意: 实际部署时应谨慎使用
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					for _, cert := range verifiedChains[0] {
						fingerprint := sha256.Sum256(cert.Raw)
						fp := hex.EncodeToString(fingerprint[:])
						if fp == client.PublicKeyFingerprint {
							return nil // 指纹匹配，证书验证成功
						}
					}
					return fmt.Errorf("TLS certificate verification failed. None of the peer certificates match the expected public key fingerprint")
				},
			},
		},
	}
}

func (client *APIClient) refreshToken() {
	// 构造获取AccessToken的请求
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

	req, err := http.NewRequest("POST", client.BaseURL+"/mgplatform/api/apps/v2/token", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating token request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		log.Printf("Error sending token request: %v", err)
		return
	}
	defer resp.Body.Close()

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("Error decoding token response: %v", err)
		return
	}

	if apiResponse.ErrNo != 0 {
		log.Printf("API error: %s (%d)", apiResponse.ErrTips, apiResponse.ErrNo)
		return
	}

	client.AccessToken = apiResponse.Data.AccessToken

	// 设置定时器以在token即将过期时自动刷新
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
}
