# SofaScore API Research & Implementation Guide

## Decision: Build Our Own SofaScore API Client ✅

Based on comprehensive research of existing Go implementations, **building our own SofaScore API client is the superior solution** compared to using football-data.org or API-Football.

## Why SofaScore?

### Advantages
1. **No API Key Required** - Public API, no registration needed
2. **No Rate Limits** - No hard 100 requests/day restriction
3. **Free Forever** - No pricing tiers or subscriptions needed
4. **Rich Data Available**:
   - Live scores with 15-second updates
   - Team standings/tables
   - Match lineups and statistics
   - Historical data
   - Head-to-head records
5. **Proven in Production** - Multiple successful Go implementations
6. **Full Control** - Custom retry logic, caching, proxy rotation
7. **Real-time Data** - Better for live scores than alternatives

### Considerations
- Unofficial API (could change without notice)
- Requires proper headers to avoid bot detection
- May need proxy support for high-volume usage
- 403 responses possible, need retry logic

## SofaScore API Overview

**Base URL**: `https://api.sofascore.com/api/v1`

### Key Endpoints for Euro21 Project

```bash
# Standings/Tables
GET /unique-tournament/{tournamentId}/season/{seasonId}/standings/total

# Scheduled Matches
GET /sport/football/scheduled-events/{date}  # Format: YYYY-MM-DD

# Live Events
GET /sport/football/events/live

# Match Details
GET /event/{eventId}
GET /event/{eventId}/lineups
GET /event/{eventId}/statistics

# Tournament Info
GET /unique-tournament/{tournamentId}/seasons
GET /category/{categoryId}/unique-tournaments
```

### Required Headers

```go
req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36")
req.Header.Set("Accept", "*/*")
req.Header.Set("Accept-Language", "en-US,en;q=0.9")
req.Header.Set("Referer", "https://www.sofascore.com/")
req.Header.Set("Sec-Fetch-Dest", "empty")
req.Header.Set("Sec-Fetch-Mode", "cors")
req.Header.Set("Sec-Fetch-Site", "same-site")
```

## Rotating Proxy: When & How

### When to Use Proxies

**Not Required For:**
- Basic CLI usage (occasional requests)
- Development and testing
- Personal use with reasonable request intervals

**Recommended For:**
- High-volume scraping (100+ requests/minute)
- Avoiding regional blocks (if you get 403 errors)
- Production applications with heavy traffic
- Distributed systems needing IP diversity

### Proxy Implementation Options

#### Option 1: SOCKS5 Proxy (Simple)
```go
import "golang.org/x/net/proxy"

func createSOCKS5Client(proxyAddr string) (*http.Client, error) {
    dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
    if err != nil {
        return nil, err
    }
    
    transport := &http.Transport{
        Dial:                dialer.Dial,
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }, nil
}
```

**Usage**:
```bash
# With SSH tunnel
ssh -D 9050 user@your-proxy-server

# In code
client, _ := createSOCKS5Client("localhost:9050")
```

#### Option 2: Tor with Circuit Rotation (Advanced)
```go
func (c *Client) RotateCircuit(ctx context.Context) error {
    // Connect to Tor control port
    conn, err := net.DialTimeout("tcp", "localhost:9051", 5*time.Second)
    if err != nil {
        return err
    }
    defer conn.Close()
    
    // Send NEWNYM signal to get new circuit
    _, err = conn.Write([]byte("AUTHENTICATE\r\n"))
    if err != nil {
        return err
    }
    
    _, err = conn.Write([]byte("SIGNAL NEWNYM\r\n"))
    if err != nil {
        return err
    }
    
    // Wait for new circuit (important!)
    time.Sleep(2 * time.Second)
    return nil
}
```

**Setup Tor**:
```bash
# Install Tor
brew install tor  # macOS
apt install tor   # Linux

# Configure /etc/tor/torrc or ~/.tor/torrc
SocksPort 9050
ControlPort 9051

# Start Tor
tor
```

**In Code**:
```go
// Create Tor client
client, _ := createSOCKS5Client("localhost:9050")

// On 403 error, rotate circuit
if resp.StatusCode == 403 {
    log.Println("403 detected, rotating Tor circuit...")
    if err := c.RotateCircuit(ctx); err != nil {
        log.Printf("Rotation failed: %v", err)
    }
    time.Sleep(2 * time.Second)
    // Retry request
}
```

### Recommended Approach for euro21

**Phase 1 (Start Simple)**:
```go
// No proxy, just proper headers
client := &http.Client{Timeout: 15 * time.Second}
```

**Phase 2 (Add Optional Proxy)**:
```go
// Support SOCKS5 proxy via config
if config.ProxyEnabled {
    client = createProxyClient(config.ProxyAddr)
}
```

**Phase 3 (If Needed: Tor Rotation)**:
```go
// Only if experiencing persistent 403s
if resp.StatusCode == 403 && config.TorEnabled {
    rotateCircuit()
    retryRequest()
}
```

## Reference Go Projects

### 1. Tactify (Most Comprehensive) ⭐⭐⭐⭐⭐
**URL**: https://github.com/imadeddine-belkat/Tactify

**Features**:
- Tor proxy with circuit rotation
- Fallback to headless Chrome (chromedp)
- Comprehensive DTO models
- Microservices architecture
- Production-ready error handling

**Key Files**:
- Client: `sofascore-service/internal/api/sofascore_api.go`
- DTOs: `shared/sofascore_models/`
- Config: `sofascore-service/config/config.go`

**Best For**: Learning advanced patterns (proxy rotation, fallbacks)

### 2. betty2310/sofascore-crawler (Simple) ⭐⭐⭐
**URL**: https://github.com/betty2310/sofascore-crawler

**Features**:
- Clean, straightforward implementation
- No proxy (proves it's not mandatory)
- Simple DTOs
- Good for beginners

**Key Files**:
- API: `api/server/event.go`
- Types: `api/types/`

**Best For**: Starting point, understanding basics

### 3. unickorn/sofacal (Proxy Example) ⭐⭐⭐
**URL**: https://github.com/unickorn/sofacal

**Features**:
- SOCKS5 proxy support
- Calendar scraping
- Minimal but functional

**Key Files**:
- Proxy setup: `event.go`

**Best For**: Learning SOCKS5 proxy integration

## Implementation Blueprint for euro21

### Project Structure
```
adapters/
  sofascore/
    client.go           # Main HTTP client
    endpoints.go        # URL builders
    proxy.go            # Optional proxy support
    cache.go            # Response caching
    retry.go            # Retry logic
    
dtos/
  sofascore/
    event.go            # Match DTOs
    standing.go         # Table DTOs
    team.go             # Team DTOs
    common.go           # Shared types
    
config/
  config.go             # Configuration
  
utils/
  ratelimiter/
    limiter.go          # Rate limiting
```

### Client Interface
```go
type FootballAPIClient interface {
    GetStandings(ctx context.Context, tournamentID, seasonID int) (*dtos.Standings, error)
    GetScheduledEvents(ctx context.Context, date string) (*dtos.Events, error)
    GetEvent(ctx context.Context, eventID int) (*dtos.EventDetail, error)
    GetLineups(ctx context.Context, eventID int) (*dtos.Lineups, error)
}
```

### Configuration
```go
type Config struct {
    // API settings
    BaseURL   string
    Timeout   time.Duration
    UserAgent string
    
    // Proxy settings (optional)
    ProxyEnabled bool
    ProxyType    string  // "socks5" or "tor"
    ProxyAddr    string  // "localhost:9050"
    
    // Tor specific (optional)
    TorControlAddr string  // "localhost:9051"
    TorEnabled     bool
    
    // Performance
    CacheEnabled  bool
    CacheTTL      time.Duration
    MaxRetries    int
    RetryDelay    time.Duration
}
```

### Sample Usage
```go
// Create client
cfg := config.LoadConfig()
client := sofascore.NewClient(cfg)

// Get Euro 2024 standings
standings, err := client.GetStandings(ctx, 1, 52186)
if err != nil {
    log.Fatalf("Failed: %v", err)
}

// Get today's matches
today := time.Now().Format("2006-01-02")
events, err := client.GetScheduledEvents(ctx, today)
```

## Key DTOs (from Research)

### Event/Match
```go
type Event struct {
    ID             int    `json:"id"`
    StartTimestamp int64  `json:"startTimestamp"`
    Slug           string `json:"slug"`
    
    Tournament struct {
        Name string `json:"name"`
        ID   int    `json:"id"`
    } `json:"tournament"`
    
    HomeTeam Team  `json:"homeTeam"`
    AwayTeam Team  `json:"awayTeam"`
    
    Status struct {
        Code        int    `json:"code"`
        Description string `json:"description"`
        Type        string `json:"type"`
    } `json:"status"`
    
    HomeScore Score `json:"homeScore"`
    AwayScore Score `json:"awayScore"`
}
```

### Standings
```go
type Standings struct {
    Standings []struct {
        Tournament struct {
            Name string `json:"name"`
            Slug string `json:"slug"`
        } `json:"tournament"`
        
        Type string `json:"type"`
        Name string `json:"name"`
        Rows []Row  `json:"rows"`
    } `json:"standings"`
}

type Row struct {
    Team struct {
        Name      string `json:"name"`
        ShortName string `json:"shortName"`
        ID        int    `json:"id"`
    } `json:"team"`
    
    Position      int `json:"position"`
    Matches       int `json:"matches"`
    Wins          int `json:"wins"`
    Draws         int `json:"draws"`
    Losses        int `json:"losses"`
    Points        int `json:"points"`
    ScoresFor     int `json:"scoresFor"`
    ScoresAgainst int `json:"scoresAgainst"`
}
```

## Best Practices

### 1. Always Handle Errors
```go
// Bad
data, _ := client.GetStandings(ctx, tid, sid)

// Good
data, err := client.GetStandings(ctx, tid, sid)
if err != nil {
    return fmt.Errorf("failed to fetch standings: %w", err)
}
```

### 2. Use Context for Timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

data, err := client.GetStandings(ctx, tid, sid)
```

### 3. Implement Caching
```go
func (c *Client) GetWithCache(endpoint string) ([]byte, error) {
    if cached := c.cache.Get(endpoint); cached != nil {
        return cached, nil
    }
    
    data, err := c.get(endpoint)
    if err == nil {
        c.cache.Set(endpoint, data, 5*time.Minute)
    }
    return data, err
}
```

### 4. Rate Limiting
```go
type RateLimiter struct {
    lastRequest time.Time
    minInterval time.Duration
    mu          sync.Mutex
}

func (r *RateLimiter) Wait() {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    elapsed := time.Since(r.lastRequest)
    if elapsed < r.minInterval {
        time.Sleep(r.minInterval - elapsed)
    }
    r.lastRequest = time.Now()
}
```

## Troubleshooting

### Problem: Getting 403 Forbidden

**Solutions**:
1. Check headers (especially User-Agent and Referer)
2. Add delay between requests (500ms+)
3. Enable proxy if available
4. Rotate Tor circuit if using Tor
5. Check if IP is blocked regionally

### Problem: Empty or Invalid JSON

**Solutions**:
1. Verify endpoint URL is correct
2. Check tournament/season IDs are valid
3. Ensure date format is YYYY-MM-DD
4. Log raw response body for debugging

### Problem: Slow Responses

**Solutions**:
1. Implement response caching
2. Use connection pooling (MaxIdleConns)
3. Reduce timeout if too high
4. Consider parallel requests for multiple data

## Testing Strategy

### Manual Testing
```bash
# Test with curl
curl -H "User-Agent: Mozilla/5.0" \
     -H "Referer: https://www.sofascore.com/" \
     "https://api.sofascore.com/api/v1/sport/football/events/live"
```

### Unit Tests
```go
func TestGetStandings(t *testing.T) {
    // Mock HTTP response
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"standings": []}`))
    }))
    defer server.Close()
    
    client := NewClient(Config{BaseURL: server.URL})
    _, err := client.GetStandings(context.Background(), 1, 1)
    assert.NoError(t, err)
}
```

## Summary

**For euro21 project**:
1. ✅ Start with simple HTTP client + proper headers
2. ✅ Add optional SOCKS5 proxy support via config
3. ✅ Implement retry logic for 403 responses
4. ✅ Add caching to reduce API load
5. ⚠️  Only add Tor rotation if experiencing persistent blocks

**Estimated complexity**: Medium
**Estimated time**: 3-4 days for core implementation
**Maintenance**: Low (stable API, proven patterns)

This approach gives us full control, zero costs, and unlimited requests while learning from proven Go implementations.
