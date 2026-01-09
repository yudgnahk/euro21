package sofascore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yudgnahk/euro21/dtos"
)

const (
	BaseURL        = "https://api.sofascore.com/api/v1"
	MaxRetries     = 3
	InitialBackoff = 1 * time.Second
)

// Client is the SofaScore API client
type Client struct {
	httpClient     *http.Client
	baseURL        string
	userAgent      string
	proxyConfig    *ProxyConfig
	cache          *Cache
	maxRetries     int
	initialBackoff time.Duration
}

// NewClient creates a new SofaScore API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		baseURL:        BaseURL,
		userAgent:      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
		maxRetries:     MaxRetries,
		initialBackoff: InitialBackoff,
	}
}

// SetProxyConfig sets the proxy configuration for the client
func (c *Client) SetProxyConfig(config *ProxyConfig) {
	c.proxyConfig = config
}

// SetCache sets the cache for the client
func (c *Client) SetCache(cache *Cache) {
	c.cache = cache
}

// newRequest creates a new HTTP request with proper headers
func (c *Client) newRequest(ctx context.Context, endpoint string) (*http.Request, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Critical headers to avoid bot detection
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://www.sofascore.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")

	return req, nil
}

// get performs a GET request and returns the response body
func (c *Client) get(ctx context.Context, endpoint string) ([]byte, error) {
	// Check cache first
	if c.cache != nil {
		if data, found := c.cache.Get(endpoint); found {
			return data, nil
		}
	}

	// Make request with retry
	data, err := c.getWithRetry(ctx, endpoint, 0)
	if err != nil {
		return nil, err
	}

	// Cache the response
	if c.cache != nil {
		if err := c.cache.Set(endpoint, data); err != nil {
			logrus.Warnf("Failed to cache response: %v", err)
		}
	}

	return data, nil
}

// getWithRetry performs a GET request with retry logic
func (c *Client) getWithRetry(ctx context.Context, endpoint string, attempt int) ([]byte, error) {
	req, err := c.newRequest(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("GET %s (attempt %d/%d)", req.URL.String(), attempt+1, c.maxRetries+1)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logrus.Warnf("Failed to close response body: %v", cerr)
		}
	}()

	// Handle success
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		return body, nil
	}

	// Handle 403 Forbidden (bot detection)
	if resp.StatusCode == http.StatusForbidden {
		if attempt < c.maxRetries {
			// Try rotating Tor circuit if configured
			if c.proxyConfig != nil && c.proxyConfig.Enabled && c.proxyConfig.Type == "tor" {
				logrus.Warnf("403 Forbidden, rotating Tor circuit...")
				if err := RotateTorCircuit(c.proxyConfig.TorControlAddr); err != nil {
					logrus.Errorf("Failed to rotate circuit: %v", err)
				}
			}

			// Exponential backoff
			backoff := c.initialBackoff * time.Duration(1<<uint(attempt))
			logrus.Infof("Retrying after %v...", backoff)
			time.Sleep(backoff)

			return c.getWithRetry(ctx, endpoint, attempt+1)
		}
		return nil, fmt.Errorf("max retries exceeded, status: 403 Forbidden")
	}

	// Handle 429 Too Many Requests
	if resp.StatusCode == http.StatusTooManyRequests {
		if attempt < c.maxRetries {
			backoff := c.initialBackoff * time.Duration(1<<uint(attempt+1)) // Longer backoff for rate limits
			logrus.Warnf("Rate limited (429), waiting %v before retry...", backoff)
			time.Sleep(backoff)
			return c.getWithRetry(ctx, endpoint, attempt+1)
		}
		return nil, fmt.Errorf("max retries exceeded, status: 429 Too Many Requests")
	}

	// Other error status codes
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// getAndUnmarshal performs a GET request and unmarshals the JSON response
func (c *Client) getAndUnmarshal(ctx context.Context, endpoint string, target interface{}) error {
	data, err := c.get(ctx, endpoint)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// GetStandings retrieves standings for a specific tournament and season
func (c *Client) GetStandings(ctx context.Context, tournamentID, seasonID int) (*dtos.SofaStandingsResponse, error) {
	endpoint := fmt.Sprintf("/unique-tournament/%d/season/%d/standings/total", tournamentID, seasonID)

	var response dtos.SofaStandingsResponse
	if err := c.getAndUnmarshal(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get standings: %w", err)
	}

	return &response, nil
}

// GetScheduledEvents retrieves scheduled events for a specific date
func (c *Client) GetScheduledEvents(ctx context.Context, date string) (*dtos.SofaEventsResponse, error) {
	endpoint := fmt.Sprintf("/sport/football/scheduled-events/%s", date)

	var response dtos.SofaEventsResponse
	if err := c.getAndUnmarshal(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get scheduled events: %w", err)
	}

	return &response, nil
}

// GetEvent retrieves details for a specific event
func (c *Client) GetEvent(ctx context.Context, eventID int) (*dtos.SofaEventDetail, error) {
	endpoint := fmt.Sprintf("/event/%d", eventID)

	var response dtos.SofaEventDetail
	if err := c.getAndUnmarshal(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &response, nil
}

// GetLiveEvents retrieves all live football events
func (c *Client) GetLiveEvents(ctx context.Context) (*dtos.SofaEventsResponse, error) {
	endpoint := "/sport/football/events/live"

	var response dtos.SofaEventsResponse
	if err := c.getAndUnmarshal(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get live events: %w", err)
	}

	return &response, nil
}

// GetTournamentEvents retrieves all events for a specific tournament and season
func (c *Client) GetTournamentEvents(ctx context.Context, tournamentID, seasonID int) (*dtos.SofaEventsResponse, error) {
	endpoint := fmt.Sprintf("/unique-tournament/%d/season/%d/events/last/0", tournamentID, seasonID)

	var response dtos.SofaEventsResponse
	if err := c.getAndUnmarshal(ctx, endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get tournament events: %w", err)
	}

	return &response, nil
}
