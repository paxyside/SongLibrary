package external_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"effictiveMobile/pkg/config"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: config.Config.ExternalApiUrl(),
	}
}

// GetSongDetails выполняет запрос к внешнему API для получения деталей о песне
func (c *Client) GetSongDetails(ctx context.Context, group, song string) (*SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", c.baseURL, group, song)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get song details: %s", resp.Status)
	}

	var detail SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, err
	}

	return &detail, nil
}
