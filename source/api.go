package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"sports_order/common"
)

// httpRequest 发送 HTTP 请求，并统一处理请求头、超时与非 200 的错误响应。
func httpRequest(method, url string, data []byte, authToken string) ([]byte, error) {
	var body io.Reader
	if data != nil {
		body = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 认证信息与标准请求头
	if authToken != "" {
		req.Header.Set("Authorization", authToken)
	}
	for key, value := range common.HTTPHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: time.Duration(common.DefaultTimeoutSec) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorBody := make([]byte, 1024)
		n, _ := resp.Body.Read(errorBody)
		if n > 0 {
			return nil, fmt.Errorf("HTTP %d: %s - %s", resp.StatusCode, resp.Status, errorBody[:n])
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// HTTPClient 是对外部 API 的简单 HTTP 实现（满足 APIClient 接口）。
type HTTPClient struct{}

// NewHTTPClient 创建一个新的 HTTPClient。
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

// Get 发起 GET 请求（不带授权）。
func (c *HTTPClient) Get(url string) ([]byte, error) {
	return httpRequest(http.MethodGet, url, nil, "")
}

// Post 发起 POST 请求（可携带授权 token）。
func (c *HTTPClient) Post(url string, data []byte, auth string) ([]byte, error) {
	return httpRequest(http.MethodPost, url, data, auth)
}
