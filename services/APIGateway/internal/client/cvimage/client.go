package cvimage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type AuthResponse struct {
	UserID         string   `json:"userId"`
	Authenticated  bool     `json:"authenticated"`
	Confidence     *float64 `json:"confidence"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) UploadScreenshot(ctx context.Context, employeeID, userID string, body io.Reader, filename string) (*http.Response, error) {
	url := fmt.Sprintf("%s/screenshot/%s", c.baseURL, employeeID)
	return c.uploadRaw(ctx, http.MethodPost, url, body, filename, map[string]string{
		"X-User-Id": userID,
	})
}

func (c *Client) GetScreenshots(ctx context.Context, employeeID, date string) (*http.Response, error) {
	url := fmt.Sprintf("%s/screenshot/%s/%s", c.baseURL, employeeID, date)
	return c.get(ctx, url)
}

func (c *Client) GetScreenshotsArchive(ctx context.Context, employeeID, date string) (*http.Response, error) {
	url := fmt.Sprintf("%s/screenshot/%s/%s/archive", c.baseURL, employeeID, date)
	return c.get(ctx, url)
}

func (c *Client) GetScreenshotFile(ctx context.Context, employeeID, filename string) (*http.Response, error) {
	url := fmt.Sprintf("%s/screenshot/%s/file/%s", c.baseURL, employeeID, filename)
	return c.get(ctx, url)
}

func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *Client) UploadPhoto(ctx context.Context, userID string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/photo", c.baseURL)
	return c.uploadRaw(ctx, http.MethodPost, url, body, userID+".png", map[string]string{
		"X-User-Id": userID,
	})
}

func (c *Client) Authenticate(ctx context.Context, body io.Reader) (*AuthResponse, int, error) {
	url := fmt.Sprintf("%s/auth", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "image/png")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, nil
	}

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp.StatusCode, err
	}
	return &result, resp.StatusCode, nil
}

func (c *Client) DeletePhoto(ctx context.Context, userID string) error {
	url := fmt.Sprintf("%s/photo/%s", c.baseURL, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("cv service returned status %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) uploadRaw(
	ctx context.Context,
	method, url string,
	body io.Reader,
	filename string,
	headers map[string]string,
) (*http.Response, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, body); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return c.httpClient.Do(req)
}
