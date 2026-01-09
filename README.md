# euro21

Euro 2021 is a modern terminal tool to get information about UEFA Euro 2021 tournament, powered by the [SofaScore API](https://www.sofascore.com/).

![Euro 2021 CLI](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

✅ **Group Stage Tables** - View standings for all 6 groups with color-coded qualification status
- 🟢 Qualified teams (top 2 + best 4 third-place teams)
- 🟡 Best third-place teams that qualify
- ⚪ Eliminated teams

✅ **Match Results** - Browse knockout stage matches (Round of 16, Quarter-finals, Semi-finals, Final)
- Match dates and times
- Team names with country flags 🇮🇹 🇫🇷 🇩🇪 🇪🇸
- Live scores and match status

✅ **Modern API Integration**
- Powered by SofaScore API (no API key required)
- File-based caching with TTL support
- Optional SOCKS5/Tor proxy support
- Retry logic with exponential backoff

## Installation

### From Source
```bash
go install github.com/yudgnahk/euro21@latest
```

### Build Locally
```bash
git clone https://github.com/yudgnahk/euro21.git
cd euro21
go build -o euro21
./euro21
```

## Usage

### Interactive Mode (Recommended)
Simply run the command and select from the menu:
```bash
euro21
```

You'll see:
```
Select menu?
  👉 Get group stage table
     Get match list
```

### View Group Stage Tables
Shows standings for all 6 groups (A-F) with complete statistics:
- Points (P), Wins (W), Draws (D), Losses (L)
- Goals For (F), Goals Against (A), Goal Difference (GD)
- Color-coded qualification indicators

### View Match List
Browse knockout stage matches with filters:
- All Matches
- Round of 16
- Quarter-finals
- Semi-finals
- Final

Displays:
- Match date and time
- Team names with country flags
- Scores (for finished matches)
- Match status (FT, LIVE, Scheduled)

## Configuration

The tool uses `euro21.config.json` for settings. Create one in the current directory or `~/.config/euro21/config.json`:

```json
{
  "proxy": {
    "enabled": false,
    "type": "socks5",
    "address": "localhost:9050",
    "tor_control_addr": "localhost:9051"
  },
  "cache": {
    "enabled": true,
    "cache_dir": "~/.cache/euro21",
    "ttl": 300
  },
  "tournament": {
    "id": 1,
    "season_id": 26542,
    "name": "UEFA Euro 2021"
  }
}
```

### Configuration Options

**Proxy Settings** (optional):
- `enabled`: Enable/disable proxy (default: false)
- `type`: Proxy type - "socks5" or "tor"
- `address`: Proxy address (e.g., "localhost:9050")
- `tor_control_addr`: Tor control port for circuit rotation (e.g., "localhost:9051")

**Cache Settings**:
- `enabled`: Enable/disable caching (default: true)
- `cache_dir`: Directory for cache files (default: "~/.cache/euro21")
- `ttl`: Cache time-to-live in seconds (default: 300 = 5 minutes)

**Tournament Settings**:
- `id`: SofaScore tournament ID (default: 1 for European Championship)
- `season_id`: SofaScore season ID (26542 for Euro 2021)
- `name`: Display name

## Development

### Project Structure
```
euro21/
├── cmd/                    # CLI commands (root, table, match)
├── adapters/sofascore/     # SofaScore API client
│   ├── client.go          # HTTP client with retry logic
│   ├── cache.go           # File-based caching
│   └── proxy.go           # SOCKS5/Tor proxy support
├── config/                 # Configuration management
├── dtos/                   # Data transfer objects
├── tablewriter/            # Terminal table rendering
├── utils/                  # Utility functions
└── tests/                  # Integration tests
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test ./... -cover

# Run specific package tests
go test ./adapters/sofascore/
go test ./config/
go test ./utils/stringutil/
```

### Linting
```bash
# Install golangci-lint
brew install golangci-lint  # macOS
# or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Test Coverage
- **adapters/sofascore**: 24.0%
- **config**: 88.5%
- **utils/stringutil**: 94.1%

## Technical Details

### API
- **Provider**: [SofaScore](https://www.sofascore.com/)
- **Base URL**: `https://api.sofascore.com/api/v1`
- **Authentication**: None required (public API)
- **Rate Limiting**: Handled via caching and retry logic

### Key Endpoints Used
- `/unique-tournament/{id}/season/{season}/standings/total` - Group standings
- `/unique-tournament/{id}/season/{season}/events/last/0` - Knockout matches

### Features
- **Retry Logic**: Automatic retries with exponential backoff for 403/429 responses
- **Caching**: MD5-based file caching with configurable TTL
- **Proxy Support**: Optional SOCKS5/Tor with circuit rotation for high-volume usage
- **Error Handling**: Comprehensive error messages with logging

## Migration from Livescore API

This project has been modernized from the legacy Livescore API to SofaScore:
- ✅ Updated to Go 1.23
- ✅ Replaced dead Livescore API with working SofaScore API
- ✅ Added caching and proxy support
- ✅ Improved error handling throughout
- ✅ Added unit tests for core functionality
- ✅ Added golangci-lint configuration
- ✅ Modernized dependencies (Cobra, Logrus, etc.)

See [MODERNIZATION_PLAN.md](MODERNIZATION_PLAN.md) for details.

## Known Limitations

- **Group Stage Matches**: The SofaScore `/events/last/0` endpoint only returns knockout stage matches for completed tournaments. Group stage match data is available through standings but not individual match results.
- **Historical Data**: As Euro 2021 is a completed tournament, data is historical and will not change.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Data provided by [SofaScore](https://www.sofascore.com/)
- Country flag emojis from [go-emoji-flags](https://github.com/yudgnahk/go-emoji-flags)
- CLI framework: [Cobra](https://github.com/spf13/cobra)
- Interactive prompts: [promptui](https://github.com/manifoldco/promptui)

## Support

If you encounter any issues or have questions:
- Open an issue on [GitHub](https://github.com/yudgnahk/euro21/issues)
- Check [SOFASCORE_API_RESEARCH.md](SOFASCORE_API_RESEARCH.md) for API details
- Review [AGENTS.md](AGENTS.md) for coding guidelines

---

Made with ⚽️ by [yudgnahk](https://github.com/yudgnahk)
