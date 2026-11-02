package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gateway/internal/config"
	"gateway/internal/model/dto"
)

type PyClient struct {
	httpClient     *http.Client
	baseURL        string
	requestTimeout time.Duration
}

func NewHTTPClient(cfg *config.Config) *PyClient {
	timeoutSec := cfg.Http.RequestTimeout
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	baseURL := strings.TrimRight(cfg.Http.AgentServerAddr, "/")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8080"
	}

	return &PyClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * time.Duration(timeoutSec),
		},
		requestTimeout: time.Second * time.Duration(timeoutSec),
	}
}

// =====================================================================================================================

func (c *PyClient) GetRequestTimeout() time.Duration {
	return c.requestTimeout
}

func (c *PyClient) SearchVectors(ctx context.Context, collectionName string, req *dto.DocsSearchRequest) (*dto.DocsSearchResponse, error) {
	if collectionName == "" {
		return nil, fmt.Errorf("collection_name is required")
	}
	if req == nil || req.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal search request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/vector/search/%s", c.baseURL, collectionName)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create search request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("call agent-server search: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read agent-server response: %w", err)
	}

	var result dto.DocsSearchResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal agent-server response: %w", err)
	}
	if resp.StatusCode != http.StatusOK || result.Code != http.StatusOK {
		msg := result.Msg
		if msg == "" {
			msg = string(respBody)
		}
		return nil, fmt.Errorf("agent-server error: status=%d code=%d msg=%s", resp.StatusCode, result.Code, msg)
	}

	return &result, nil
}

func (c *PyClient) ListCollections(ctx context.Context) ([]*dto.Collection, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/vector/collections", nil)
	if err != nil {
		return nil, fmt.Errorf("create collection request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call agent-server collections: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read agent-server collections: %w", err)
	}
	var result struct {
		Code int               `json:"code"`
		Msg  string            `json:"msg"`
		Data []*dto.Collection `json:"data"`
	}
	if len(body) == 0 || json.Unmarshal(body, &result) != nil || resp.StatusCode != http.StatusOK || result.Code != http.StatusOK {
		return nil, fmt.Errorf("agent-server collections unavailable")
	}
	return result.Data, nil
}

func (c *PyClient) CreateCollection(ctx context.Context, collectionName string) error {
	var result dto.Collection
	if err := c.doJSON(ctx, http.MethodPost, "/api/v1/vector/collections/"+collectionName, nil, &result); err != nil {
		return err
	}
	return nil
}

func (c *PyClient) DeleteCollection(ctx context.Context, collectionName string) error {
	return c.doJSON(ctx, http.MethodDelete, "/api/v1/vector/collections/"+collectionName, nil, nil)
}

func (c *PyClient) ListDocuments(ctx context.Context, collectionName string, page, pageSize int) (*dto.ListDocumentsResponse, error) {
	path := fmt.Sprintf("/api/v1/vector/collections/%s/documents?page=%d&page_size=%d", collectionName, page, pageSize)
	var result dto.ListDocumentsResponse
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *PyClient) doJSON(ctx context.Context, method, path string, payload any, target any) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("call agent-server: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result := struct {
		Code int             `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}{}
	if err := json.Unmarshal(data, &result); err != nil || resp.StatusCode != http.StatusOK || result.Code != http.StatusOK {
		return fmt.Errorf("agent-server request failed: %s", result.Msg)
	}
	if target != nil && len(result.Data) > 0 && string(result.Data) != "null" {
		return json.Unmarshal(result.Data, target)
	}
	return nil
}
