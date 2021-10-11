package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) TestAccessToken(accesstoken string) error {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://asia-northeast1-hashira-web.cloudfunctions.net/test-access-token", nil)
	if err != nil {
		return fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", accesstoken))

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("This accesstoken is not valid. Please check HASHIRA_ACCESS_TOKEN is correct and try again [%d].", resp.StatusCode)
	}

	return nil
}

type Task struct {
	ID        string
	Name      string
	Place     string
	IsDeleted bool
}

type Priority map[string][]string

// priority's key should be one of following strings:
// "BACKLOG", "TODO", "DOING", "DONE"
type UploadRequest struct {
	Tasks    map[string]Task `json:"tasks"`
	Priority Priority        `json:"priority"`
}

func (c *Client) Upload(accesstoken string, ur UploadRequest) error {
	body, err := json.Marshal(ur)
	if err != nil {
		return fmt.Errorf("failed to marshal upload request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://asia-northeast1-hashira-web.cloudfunctions.net/upload",
		bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create new request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", accesstoken))

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

type DownloadResult UploadRequest

func (c *Client) Download(accesstoken string) (DownloadResult, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://asia-northeast1-hashira-web.cloudfunctions.net/download", nil)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to prepare request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", accesstoken))

	httpcli := http.Client{}
	resp, err := httpcli.Do(req)
	if err != nil {
		return DownloadResult{}, fmt.Errorf("failed to download tasks and priorities: %w", err)
	}
	defer resp.Body.Close()

	var ret DownloadResult
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return DownloadResult{}, fmt.Errorf("failed to decode response body: %w", err)
	}

	return ret, nil
}
