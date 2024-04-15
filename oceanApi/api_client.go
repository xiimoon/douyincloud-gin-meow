// oceanApi 包用于向 Ocean Engine 的 API 发送请求。
package oceanApi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	BaseURL    string
}

// RequestData 定义了发送到 Ocean Engine 的 JSON 结构。
type RequestData struct {
	EventType string  `json:"event_type"`
	Context   Context `json:"context"`
	Timestamp int64   `json:"timestamp"`
}

// Context 定义了请求的 Context 部分。
type Context struct {
	Ad AdInfo `json:"ad"`
}

// AdInfo 定义了广告相关信息。
type AdInfo struct {
	Callback string `json:"callback"`
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		BaseURL:    "http://analytics.oceanengine.com/api/v2",
	}
}

// SendConversion 发送转换事件到 Ocean Engine。
func (c *Client) SendConversion(ctx context.Context, eventType string, clickid string) (string, error) {
	data := RequestData{
		EventType: eventType,
		Context: Context{
			Ad: AdInfo{
				Callback: clickid,
			},
		},
		Timestamp: time.Now().Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/conversion", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
