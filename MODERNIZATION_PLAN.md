# euro21 Modernization Plan

## Executive Summary

The euro21 project has multiple critical issues that need to be addressed:
- **Dead API**: Livescore API returns 503 Service Unavailable
- **Outdated Go version**: Using Go 1.16 (released Feb 2021, now EOL)
- **Deprecated packages**: Using `io/ioutil` (deprecated since Go 1.16)
- **Outdated dependencies**: All major dependencies have newer versions
- **Poor error handling**: Errors are frequently ignored
- **No tests**: Zero test coverage
- **Tournament locked to Euro 2021**: Event ended in 2021

## Critical Issues (Must Fix)

### 1. Dead Livescore API ❌ CRITICAL
**Problem**: The API endpoint `https://prod-public-api.livescore.com/v1/api/react` returns 503 errors.

**Solution: Build Custom SofaScore API Client** ✅ SELECTED

#### Why SofaScore?
- **Pros**: 
  - No API key required (public API)
  - No hard rate limits
  - Free forever
  - Rich data (live scores, lineups, statistics, standings)
  - Multiple successful Go implementations exist
  - Real-time data with 15-second updates
- **Cons**: 
  - Unofficial API (could change)
  - May require headers to avoid bot detection
  - Optional proxy support for high-volume usage
- **Cost**: FREE (no registration needed)
- **URL**: `https://api.sofascore.com/api/v1`

#### SofaScore Key Endpoints:
```
# Events/Matches
GET /sport/football/scheduled-events/{date}          # Format: YYYY-MM-DD
GET /event/{eventId}                                 # Event details
GET /event/{eventId}/lineups                         # Match lineups
GET /event/{eventId}/statistics                      # Match statistics

# Standings/Tables  
GET /unique-tournament/{tournamentId}/season/{seasonId}/standings/total

# Tournaments
GET /unique-tournament/{tournamentId}/seasons        # All seasons
GET /category/{categoryId}/unique-tournaments        # Leagues by country
```

#### Proven Go Implementations:
1. **Tactify** - Advanced with Tor proxy rotation and circuit management
2. **betty2310/sofascore-crawler** - Simple, clean implementation
3. **unickorn/sofacal** - SOCKS5 proxy support

**Implementation Steps**:
1. Create new adapter `adapters/sofascore/` with client structure
2. Implement HTTP client with proper headers (User-Agent, Referer, Sec-Fetch-*)
3. Design DTOs matching SofaScore JSON responses
4. Add optional SOCKS5/Tor proxy support with rotation
5. Implement retry logic for 403 responses (bot detection)
6. Add response caching to reduce load
7. Add rate limiting (optional, for politeness)
8. Implement context-aware requests with timeouts

### 2. Outdated Go Version ⚠️ HIGH PRIORITY
**Current**: Go 1.16 (released Feb 2021, EOL)
**Target**: Go 1.23+ (current stable)

**Changes Required**:
- Update `go.mod`: `go 1.23`
- Update `.github/workflows/go.yml`: Use `go-version: '1.23'`
- Replace deprecated `io/ioutil` with `io` and `os`

**Breaking Changes**: None expected (Go maintains backwards compatibility)

### 3. Deprecated io/ioutil Package ⚠️ HIGH PRIORITY
**Problem**: `io/ioutil` deprecated since Go 1.16

**Replacements**:
```go
// Old
body, err := ioutil.ReadAll(res.Body)

// New
body, err := io.ReadAll(res.Body)
```

**Files to Update**:
- `adapters/livescore.go:52`

### 4. Outdated Dependencies ⚠️ HIGH PRIORITY
**Current → Latest**:
- `github.com/spf13/cobra`: v1.1.3 → v1.10.2
- `github.com/sirupsen/logrus`: v1.8.1 → v1.9.3
- `github.com/manifoldco/promptui`: v0.8.0 → v0.9.0

**Update Command**:
```bash
go get -u github.com/spf13/cobra@latest
go get -u github.com/sirupsen/logrus@latest
go get -u github.com/manifoldco/promptui@latest
go mod tidy
```

**Risk**: Low - these are patch/minor version updates with backwards compatibility

## Important Issues (Should Fix)

### 5. Poor Error Handling 📝 MEDIUM PRIORITY
**Problem**: Errors are ignored throughout the codebase

**Examples**:
```go
// cmd/table.go:23
data, _ := adapters.GetTables()  // Error ignored!

// cmd/match.go:97
data, _ := adapters.GetStage(...)  // Error ignored!

// adapters/livescore.go:66
request, _ := newGetRequest(...)  // Error ignored!
```

**Solution**:
```go
data, err := adapters.GetTables()
if err != nil {
    logrus.Errorf("Failed to fetch tables: %v", err)
    return err
}
```

**Files to Fix**:
- `cmd/table.go`
- `cmd/match.go`
- `cmd/root.go`
- `adapters/livescore.go`

### 6. No Test Coverage 📝 MEDIUM PRIORITY
**Current**: 0% test coverage
**Target**: 60%+ coverage

**Priority Test Areas**:
1. `utils/stringutil` - Pure functions, easy to test
2. `utils/sliceutil` - Pure functions
3. `adapters/` - Mock HTTP responses
4. `tablewriter/` - Test rendering logic

**Example Test**:
```go
// utils/stringutil/string_test.go
func TestToInt(t *testing.T) {
    tests := []struct {
        input    string
        expected int
    }{
        {"123", 123},
        {"0", 0},
        {"invalid", 0},
    }
    
    for _, tt := range tests {
        got := ToInt(tt.input)
        if got != tt.expected {
            t.Errorf("ToInt(%q) = %d, want %d", tt.input, got, tt.expected)
        }
    }
}
```

### 7. Tournament Locked to Euro 2021 📝 MEDIUM PRIORITY
**Problem**: Hardcoded to Euro 2020/2021

**Solution Options**:
- **Option A**: Update to Euro 2024 (latest tournament)
- **Option B**: Make tournament-agnostic with flag: `euro21 --tournament euro2024 table`
- **Option C**: Rename project and support multiple tournaments

**Recommendation**: Option B - Add tournament selection while maintaining backwards compatibility

### 8. Configuration Management 📝 MEDIUM PRIORITY
**Problem**: No configuration file for proxies, settings

**Solution**: Add config file support
```go
// config/config.go
type Config struct {
    ProxyEnabled bool   `json:"proxy_enabled"`
    ProxyType    string `json:"proxy_type"`    // "socks5" or "tor"
    ProxyAddr    string `json:"proxy_addr"`    // "localhost:9050"
    TorControl   string `json:"tor_control"`   // "localhost:9051" for circuit rotation
    CacheEnabled bool   `json:"cache_enabled"`
    CacheDir     string `json:"cache_dir"`
    CacheTTL     int    `json:"cache_ttl"`     // seconds
    Tournament   string `json:"tournament"`     // "euro2024"
    UserAgent    string `json:"user_agent"`
}
```

**Config File Locations**:
1. `~/.config/euro21/config.json`
2. `./euro21.config.json`
3. Environment variables: `EURO21_PROXY_ENABLED`, `EURO21_PROXY_ADDR`

## Nice-to-Have Improvements

### 9. CI/CD Improvements 🔧 LOW PRIORITY
**Changes**:
- Update GitHub Actions to Go 1.23
- Add `go vet` and `golangci-lint` checks
- Add test coverage reporting
- Add build artifacts upload

### 10. Linting Configuration 🔧 LOW PRIORITY
**Add `.golangci.yml`**:
```yaml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
```

## Implementation Roadmap

### Phase 1: Critical Fixes (Week 1) - 3 days
1. ✅ Update Go version to 1.23
2. ✅ Replace `io/ioutil` with `io`/`os`
3. ✅ Update all dependencies (cobra, logrus, promptui)
4. ✅ Build SofaScore API client with proper structure
5. ✅ Implement core endpoints (standings, matches, events)
6. ✅ Add proper headers to avoid bot detection
7. ✅ Test basic functionality (table, match commands)

### Phase 2: Proxy & Reliability (Week 1-2) - 3 days
1. ✅ Implement SOCKS5 proxy support
2. ✅ Add optional Tor integration with circuit rotation
3. ✅ Implement retry logic for 403/429 responses
4. ✅ Add response caching (file-based or memory)
5. ✅ Add rate limiting (politeness delays)
6. ✅ Fix all error handling throughout codebase
7. ✅ Add structured logging (logrus)

### Phase 3: Features & Quality (Week 2) - 4 days
1. ✅ Add configuration file support
2. ✅ Add tournament selection (Euro 2024, Euro 2020, etc.)
3. ✅ Add unit tests (target 60% coverage)
4. ✅ Update GitHub Actions workflow
5. ✅ Add golangci-lint configuration
6. ✅ Update documentation (README, AGENTS.md)
7. ✅ Add example configurations

## Migration Strategy

### Backwards Compatibility
- Keep existing command structure: `euro21 table`, `euro21 match`
- Add optional flags without breaking existing usage
- Default to latest Euro tournament

### SofaScore API Client Architecture
```go
// Package structure
adapters/
  sofascore/
    client.go       # Main HTTP client with headers
    proxy.go        # SOCKS5/Tor proxy support
    endpoints.go    # Endpoint URL builders
    cache.go        # Response caching
    retry.go        # Retry logic for 403s
    
dtos/
  sofascore/
    event.go        # Match/event DTOs
    standing.go     # Table/standings DTOs
    team.go         # Team DTOs
    common.go       # Shared types

config/
  config.go         # Configuration management

// Client interface for future extensibility
type FootballAPIClient interface {
    GetStandings(ctx context.Context, tournamentID, seasonID int) (*dtos.Standings, error)
    GetScheduledEvents(ctx context.Context, date string) (*dtos.Events, error)
    GetEvent(ctx context.Context, eventID int) (*dtos.EventDetail, error)
    GetLineups(ctx context.Context, eventID int) (*dtos.Lineups, error)
}

// SofaScore implementation
type SofascoreClient struct {
    httpClient  *http.Client
    config      *config.Config
    cache       *Cache
    rateLimiter *RateLimiter
}
```

### Key Implementation Details

**1. HTTP Client with Headers**
```go
func (c *SofascoreClient) newRequest(ctx context.Context, endpoint string) (*http.Request, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    
    // Critical headers to avoid bot detection
    req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36")
    req.Header.Set("Accept", "*/*")
    req.Header.Set("Accept-Language", "en-US,en;q=0.9")
    req.Header.Set("Referer", "https://www.sofascore.com/")
    req.Header.Set("Sec-Fetch-Dest", "empty")
    req.Header.Set("Sec-Fetch-Mode", "cors")
    req.Header.Set("Sec-Fetch-Site", "same-site")
    
    return req, nil
}
```

**2. Proxy Support (Optional)**
```go
func createProxyClient(cfg *config.Config) (*http.Client, error) {
    if !cfg.ProxyEnabled {
        return &http.Client{Timeout: 15 * time.Second}, nil
    }
    
    dialer, err := proxy.SOCKS5("tcp", cfg.ProxyAddr, nil, proxy.Direct)
    if err != nil {
        return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
    }
    
    transport := &http.Transport{
        Dial:                dialer.Dial,
        MaxIdleConns:        10,
        IdleConnTimeout:     30 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
    }
    
    return &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }, nil
}
```

**3. Retry Logic**
```go
func (c *SofascoreClient) getWithRetry(ctx context.Context, endpoint string, maxRetries int) ([]byte, error) {
    for attempt := 0; attempt <= maxRetries; attempt++ {
        resp, err := c.httpClient.Do(req)
        
        if resp.StatusCode == http.StatusOK {
            return io.ReadAll(resp.Body)
        }
        
        if resp.StatusCode == http.StatusForbidden && c.config.ProxyEnabled {
            logrus.Warn("403 detected, rotating proxy...")
            if err := c.rotateCircuit(ctx); err != nil {
                logrus.Errorf("Circuit rotation failed: %v", err)
            }
            time.Sleep(2 * time.Second)
            continue
        }
        
        return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
    }
    
    return nil, fmt.Errorf("max retries exceeded")
}
```

## Testing Strategy

### Manual Testing Checklist
- [ ] `euro21` - Interactive menu works
- [ ] `euro21 table` - Displays group tables
- [ ] `euro21 match` - Displays match fixtures
- [ ] Error handling - Graceful failures with clear messages
- [ ] API key missing - Clear error message
- [ ] Rate limiting - Handles 429 responses
- [ ] Caching - Reduces API calls

### Automated Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./utils/stringutil

# Run single test
go test ./utils/stringutil -run TestToInt
```

## Risk Assessment

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| New API has different data structure | High | High | Thorough DTO mapping and testing |
| API rate limits hit during testing | Medium | Medium | Implement caching and mocking |
| Breaking changes in dependencies | Low | Low | Test thoroughly before release |
| Go 1.23 compatibility issues | Low | Very Low | Go maintains strong backwards compatibility |

## Success Criteria

1. ✅ Application builds without errors on Go 1.23
2. ✅ All API calls work with new provider
3. ✅ No ignored errors in codebase
4. ✅ Test coverage ≥ 60%
5. ✅ CI/CD pipeline passes
6. ✅ Documentation updated

## Estimated Effort

- **Phase 1 (Critical - SofaScore Client)**: 3 days
- **Phase 2 (Proxy & Reliability)**: 3 days  
- **Phase 3 (Features & Testing)**: 4 days
- **Total**: 10 days (2 weeks)

## Next Steps

1. ✅ Review and approve this plan
2. Create feature branch: `feature/modernization-sofascore`
3. Implement Phase 1: Core SofaScore client
   - Update Go version and dependencies
   - Build basic client with proper headers
   - Implement DTOs for events and standings
   - Test endpoints manually
4. Implement Phase 2: Proxy support
   - Add SOCKS5 proxy configuration
   - Optional: Tor with circuit rotation
   - Implement retry and caching
5. Implement Phase 3: Polish
   - Add tests
   - Update documentation
   - Add example configs
6. Test thoroughly with different tournaments
7. Merge and release v2.0.0

## Reference Projects

**Study these Go implementations:**
1. **Tactify** (Advanced): https://github.com/imadeddine-belkat/Tactify
   - Best reference for proxy rotation and circuit management
   - Comprehensive DTO models
   
2. **betty2310/sofascore-crawler** (Simple): https://github.com/betty2310/sofascore-crawler
   - Clean, straightforward implementation
   - Good starting point for basic client

3. **unickorn/sofacal** (Proxy): https://github.com/unickorn/sofacal
   - SOCKS5 proxy implementation example
