package musicbrainz

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Client handles rate-limited queries to the MusicBrainz API.
type Client struct {
	baseURL    string
	userAgent  string
	httpClient *http.Client
	mu         sync.Mutex
	lastReq    time.Time
	minDelay   time.Duration
}

// Config specifies required identification for MusicBrainz API guidelines.
type Config struct {
	AppName     string
	Version     string
	ContactInfo string
}

// NewClient creates a rate-limited MusicBrainz API client complying with 1 req/sec limits.
func NewClient(cfg Config) *Client {
	if cfg.AppName == "" {
		cfg.AppName = "IndianMusicDB"
	}
	if cfg.Version == "" {
		cfg.Version = "0.1.0"
	}
	if cfg.ContactInfo == "" {
		cfg.ContactInfo = "admin@example.com"
	}

	ua := fmt.Sprintf("%s/%s ( %s )", cfg.AppName, cfg.Version, cfg.ContactInfo)

	return &Client{
		baseURL:   "https://musicbrainz.org/ws/2",
		userAgent: ua,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		minDelay: 1 * time.Second, // 1 req/sec compliance
	}
}

// rateLimit blocks until at least 1 second has elapsed since the last request.
func (c *Client) rateLimit(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	elapsed := time.Since(c.lastReq)
	if elapsed < c.minDelay {
		wait := c.minDelay - elapsed
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
	c.lastReq = time.Now()
	return nil
}

// Do executes a rate-limited HTTP GET request with proper User-Agent header.
func (c *Client) Do(ctx context.Context, endpoint string, params url.Values, target interface{}) error {
	if err := c.rateLimit(ctx); err != nil {
		return err
	}

	if params == nil {
		params = make(url.Values)
	}
	params.Set("fmt", "json")

	reqURL := fmt.Sprintf("%s/%s?%s", c.baseURL, endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create musicbrainz request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("musicbrainz HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("musicbrainz API returned status code %d for %s", resp.StatusCode, reqURL)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode musicbrainz response: %w", err)
	}

	return nil
}

// ArtistResponse represents a simplified MusicBrainz artist query response.
type ArtistResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SortName string `json:"sort-name"`
	Type     string `json:"type"`
	Disambig string `json:"disambiguation"`
	Country  string `json:"country"`
}

// ArtistSearchResult represents the search response from MusicBrainz.
type ArtistSearchResult struct {
	Created string           `json:"created"`
	Count   int              `json:"count"`
	Offset  int              `json:"offset"`
	Artists []ArtistResponse `json:"artists"`
}

// SearchArtist queries MusicBrainz for artists by name.
func (c *Client) SearchArtist(ctx context.Context, query string) (*ArtistSearchResult, error) {
	var resp ArtistSearchResult
	params := url.Values{}
	params.Set("query", fmt.Sprintf("artist:%s", query))

	if err := c.Do(ctx, "artist", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
