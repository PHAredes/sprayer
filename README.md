# Job Scraper - TUI Application

A powerful job scraping application with a beautiful terminal UI built with Go and Charm libraries. Scrapes jobs from multiple sources including LinkedIn, Indeed, Glassdoor, Reddit, RSS feeds, and APIs.

## Features

- **Beautiful TUI Interface**: Built with Charm's Bubble Tea framework
- **Multi-Source Scraping**: Python and Go-based scrapers for maximum coverage
- **Proxy Support**: Built-in proxy rotation to avoid IP blocks
- **Real-time Progress**: Live updates during scraping
- **Smart Filtering**: Filter by keywords, score, location, and company
- **Export Options**: Save results to JSON format
- **Score-based Ranking**: Jobs scored based on CV keyword matching
- **TDD Approach**: Comprehensive test coverage

## Installation

### Prerequisites

- Go 1.21 or later
- Python 3.8 or later
- Docker (optional, for containerized execution)

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd job-scraper
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install Python dependencies:
```bash
pip install -r requirements.txt
```

4. (Optional) Set up proxies:
```bash
# Edit config/proxies.txt with your proxy list
cp config/proxies.txt.example config/proxies.txt
```

## Usage

### Running the Application

```bash
# Run the TUI application
go run main_with_scrapers.go
```

### Key Controls

- `s` - Start scraping jobs
- `f` - Open filter menu
- `Enter` - View job details
- `Esc` - Go back / Stop scraping
- `r` - Refresh job list
- `e` - Export jobs
- `q` - Quit

### Python Scrapers

The system includes Python scrapers for specialized sources:

```bash
# Run enhanced scraper with proxy support
python scripts/enhanced_scraper_with_proxy.py

# Run Greenhouse scraper
python scripts/greenhouse_scraper.py
```

## Architecture

### Components

1. **Go TUI Application** (`main_with_scrapers.go`)
   - Terminal user interface
   - Scraping orchestration
   - Real-time updates

2. **Go Scrapers** (`scrapers/scrapers.go`)
   - Web scraping with Colly
   - API integration
   - RSS feed parsing

3. **Python Scrapers** (`scripts/`)
   - LinkedIn scraper with proxy support
   - Greenhouse job board scraper
   - Enhanced scraper with anti-scraping techniques

4. **Configuration** (`config/`)
   - Proxy settings
   - API keys
   - Scraper settings

### Job Scoring System

Jobs are scored based on keyword matching:

- **High-value keywords** (10 points): Rust, C++, Haskell, LLVM, Compiler, WebAssembly
- **Medium-value keywords** (7-8 points): Go, Embedded, Systems, Performance, Algorithms
- **Low-value keywords** (3-5 points): Git, Docker, Remote, Developer

Blocked keywords reduce the score by 10 points each.

## Testing

### Run All Tests

```bash
# Run Go tests
go test ./...

# Run with coverage
go test -cover ./...

# Run Python tests
python -m pytest tests/
```

### Test Coverage

- Unit tests for all components
- Integration tests for full workflow
- Benchmark tests for performance
- TDD approach with comprehensive coverage

## Configuration

### Proxy Configuration

Edit `config/proxies.txt`:
```
# Add your proxies here
http://proxy1.example.com:8080
http://user:pass@proxy2.example.com:8080
socks5://proxy3.example.com:1080
```

### Scraper Settings

Edit `config/scraper_config.json`:
```json
{
  "max_workers": 20,
  "rate_limit_delay": 2.0,
  "user_agents": [...],
  "blocked_keywords": [...],
  "search_queries": [...]
}
```

## Output

Jobs are saved to `/tmp/job-scraper/` (or your configured output directory):
- `jobs_[scraper]_[timestamp].json` - Individual scraper results
- `jobs_combined_[timestamp].json` - All jobs combined
- `jobs_export_[timestamp].json` - Exported filtered jobs

## Development

### TDD Approach

1. Write tests first
2. Implement functionality
3. Refactor and optimize
4. Ensure all tests pass

### Adding New Scrapers

1. Implement the `Scraper` interface:
```go
type Scraper interface {
    Scrape() ([]Job, error)
    GetName() string
}
```

2. Add tests in `*_test.go`
3. Register in `CreateDefaultScrapers()`
4. Update documentation

### Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for complex logic
- Keep functions small and focused

## Performance

- Concurrent scraping with goroutines
- Rate limiting to avoid blocks
- Efficient data structures
- Minimal memory footprint

## Troubleshooting

### Common Issues

1. **Proxy not working**: Check proxy format and connectivity
2. **Scrapers blocked**: Increase delays and rotate proxies
3. **No jobs found**: Check filters and keywords
4. **TUI not rendering**: Check terminal compatibility

### Debug Mode

Enable debug logging:
```bash
export DEBUG=true
go run main_with_scrapers.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for beautiful TUI libraries
- [Colly](https://github.com/gocolly/colly) for web scraping
- [Scrapy](https://scrapy.org/) for inspiration
- The open-source community for various scraping techniques

## FAQ

### Q: How many jobs can it find?
A: The system is designed to find 300+ jobs from multiple sources. Actual count depends on your filters and market conditions.

### Q: Is it legal to scrape job sites?
A: Scraping for personal use is generally acceptable, but respect robots.txt and rate limits. For commercial use, check terms of service.

### Q: How do I avoid getting blocked?
A: Use rotating proxies, respect rate limits, and vary user agents. The system includes built-in protection.

### Q: Can I add custom job boards?
A: Yes! Implement the Scraper interface and add it to the manager. See the documentation for details.