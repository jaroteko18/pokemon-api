package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SupabaseClient handles direct HTTP requests to Supabase REST API
type SupabaseClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewSupabaseClient creates a new Supabase HTTP client
func NewSupabaseClient(url, apiKey string) *SupabaseClient {
	return &SupabaseClient{
		baseURL:    url + "/rest/v1",
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Insert performs a POST request to insert data
func (c *SupabaseClient) Insert(table string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	url := fmt.Sprintf("%s/%s", c.baseURL, table)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// SelectAllPaginated performs a GET request with pagination
func (c *SupabaseClient) SelectAllPaginated(table string, page, limit int) ([]byte, error) {
	offset := (page - 1) * limit
	url := fmt.Sprintf("%s/%s?select=*&order=registered_at.desc&limit=%d&offset=%d", c.baseURL, table, limit, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Prefer", "count=exact")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// SelectAll performs a GET request to fetch all records
func (c *SupabaseClient) SelectAll(table string) ([]byte, error) {
	return c.SelectAllOrdered(table, "registered_at.desc")
}

// SelectAllOrdered performs a GET request to fetch all records with custom ordering
func (c *SupabaseClient) SelectAllOrdered(table string, order string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?select=*&order=%s", c.baseURL, table, order)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Select performs a GET request with optional filters
func (c *SupabaseClient) Select(table, column, value string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s?%s=eq.%s&select=*", c.baseURL, table, column, value)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Update performs a PATCH request
func (c *SupabaseClient) Update(table, column, value string, updates map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(updates)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal updates: %w", err)
	}

	url := fmt.Sprintf("%s/%s?%s=eq.%s", c.baseURL, table, column, value)
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", c.apiKey)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Prefer", "return=representation")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}
