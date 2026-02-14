# Sprayer

Sprayer is an **agentic job application tool** designed to automate the painful parts of job hunting. It combines high-performance scraping, LLM-based personalization, and a terminal-native workflow.

## Features

- **Multi-Source Scraping**: Scrapes Hacker News (Who is Hiring), RemoteOK, WeWorkRemotely, Arbeitnow, Jobicy, LinkedIn, Indeed, Glassdoor, and specialized RSS feeds (Golang, Rust, etc.).
- **LLM Integration**: Uses OpenAI-compatible APIs (e.g. iFlow / Moonshot K2) to generate personalized cover letters and emails.
- **TUI & CLI**: Beautiful terminal user interface (Bubble Tea) and scriptable CLI.
- **Compositional Design**: Unix-philosophy architecture — scrapers, filters, and matchers are composable pipelines.
- **Local-First**: fast SQLite storage, local profiles, mu4e-compatible email drafts.

## Installation

```bash
git clone https://github.com/user/sprayer.git
cd sprayer
go mod tidy
go build -o sprayer ./cmd/sprayer/main.go
go build -o sprayer-cli ./cmd/cli/main.go
```

## Configuration

Set up your LLM credentials (required for email generation):

```bash
export SPRAYER_LLM_URL="https://apis.iflow.cn/v1"  # or https://api.openai.com/v1
export SPRAYER_LLM_KEY="your-api-key"
export SPRAYER_LLM_MODEL="kimi-k2"                 # or gpt-4o, deepseek-v3, etc.
```

## Usage

### Interactive TUI

Run `./sprayer` to enter the interactive mode.

- **s**: Scrape new jobs
- **f**: Filter by keywords
- **p**: Switch profiles
- **a**: Apply (generate email draft)
- **j/k**: Navigation
- **Enter**: View details

### CLI Automation

Scrape and save to DB:
```bash
./sprayer-cli scrape "rust" "remote"
```

List and filter jobs:
```bash
./sprayer-cli list --keywords "rust,compiler" --min-score 80
```

Apply to a specific job (generates draft):
```bash
./sprayer-cli apply --job "hn-123456" --prompt "email_cold"
```

## Project Structure

- `cmd/`: Entrypoints (`sprayer`, `cli`)
- `internal/job/`: Core domain model, storage, and pipeline logic
- `internal/scraper/`: Concrete scrapers (HN, LinkedIn, etc.)
- `internal/profile/`: Profile management and matching logic
- `internal/llm/`: LLM client and prompt management
- `internal/apply/`: Email generation and export
- `internal/ui/`: TUI and CLI implementation
- `prompts/`: Text templates for LLM generation

## License

MIT