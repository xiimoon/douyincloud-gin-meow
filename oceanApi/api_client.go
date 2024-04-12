package oceanApi

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	BaseURL string
	AppID   string
	Secret  string

	AccessToken string
}

type APIResponse struct {
	Data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	} `json:"data"`
}

func NewAPIClient(baseURL, appID, secret string) {

	url := "https://minigame.zijieapi.com/mgplatform/api/apps/v2/token"

	// 定义请求的数据
	data := fmt.Sprintf(`{"appid":"%s","secret":"%s","grant_type":"client_credential"}`, appID, secret)

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 创建一个 HTTP 客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// 打印响应体
	fmt.Println(string(body))
}
