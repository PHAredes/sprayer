#!/bin/bash

# Build script for Job Scraper TUI application

set -e

echo "🚀 Building Job Scraper TUI Application..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_status "Go version: $GO_VERSION"
}

# Check if Python is installed
check_python() {
    if ! command -v python3 &> /dev/null; then
        print_error "Python 3 is not installed. Please install Python 3.8 or later."
        exit 1
    fi
    
    PYTHON_VERSION=$(python3 --version | awk '{print $2}')
    print_status "Python version: $PYTHON_VERSION"
}

# Install Go dependencies
install_go_deps() {
    print_status "Installing Go dependencies..."
    go mod download
    go mod tidy
}

# Install Python dependencies
install_python_deps() {
    print_status "Installing Python dependencies..."
    
    # Create virtual environment if it doesn't exist
    if [ ! -d "venv" ]; then
        print_status "Creating Python virtual environment..."
        python3 -m venv venv
    fi
    
    # Activate virtual environment
    source venv/bin/activate
    
    # Create requirements.txt if it doesn't exist
    if [ ! -f "requirements.txt" ]; then
        print_warning "requirements.txt not found, creating it..."
        cat > requirements.txt << EOF
requests>=2.31.0
beautifulsoup4>=4.12.0
feedparser>=6.0.10
selenium>=4.15.0
lxml>=4.9.3
EOF
    fi
    
    # Install dependencies
    pip install -r requirements.txt
    
    # Deactivate virtual environment
    deactivate
}

# Create necessary directories
create_dirs() {
    print_status "Creating necessary directories..."
    mkdir -p config
    mkdir -p outputs
    mkdir -p scripts
    mkdir -p logs
}

# Create default config files
create_configs() {
    print_status "Creating default configuration files..."
    
    # Create proxy config if it doesn't exist
    if [ ! -f "config/proxies.txt" ]; then
        cat > config/proxies.txt << EOF
# Proxy List for Enhanced Job Scraper
# Add your proxies here, one per line
# Format options:
# - HTTP: http://ip:port
# - HTTP with auth: http://username:password@ip:port
# - SOCKS: socks5://ip:port
# - SOCKS with auth: socks5://username:password@ip:port

# Free Proxies (for testing - these may not work well)
# http://192.168.1.1:8080
# http://proxy.example.com:3128
# socks5://127.0.0.1:1080

# Paid Proxy Services (recommended for production):
# - Bright Data (formerly Luminati)
# - Oxylabs
# - Smartproxy
# - Storm Proxies
# - Blazing SEO

# Note: For LinkedIn scraping, use residential proxies for best results
EOF
    fi
    
    # Create scraper config
    if [ ! -f "config/scraper_config.json" ]; then
        cat > config/scraper_config.json << EOF
{
  "max_workers": 20,
  "rate_limit_delay": 2.0,
  "user_agents": [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
  ],
  "blocked_keywords": [
    "outsource", "microempresa", "erp", "consultoria operacional",
    "infraestrutura legacy", "delphi", "cobol", "ui/ux designer",
    "mobile", "lead", "tech lead", "estágio", "estagiario",
    "trainee", "presencial", "front", "frontend", "front-end",
    "python", "php", "dba", "redator", "vendas", "gerente",
    "salesforce", "portugal", "servicenow", "soap", "vue",
    "nuxt", "nest", "7 years", "bellroy", "intern", "junior",
    "entry level", "graduate", "apprentice", "senior", "principal"
  ],
  "search_queries": [
    {"keywords": "rust developer", "location": "Remote"},
    {"keywords": "haskell developer", "location": "Remote"},
    {"keywords": "c++ developer", "location": "Remote"},
    {"keywords": "go developer", "location": "Remote"},
    {"keywords": "systems programmer", "location": "Remote"},
    {"keywords": "llvm developer", "location": "Remote"},
    {"keywords": "embedded engineer", "location": "Remote"},
    {"keywords": "compiler engineer", "location": "Remote"}
  ]
}
EOF
    fi
}

# Run tests
run_tests() {
    print_status "Running tests..."
    
    # Run Go tests
    print_status "Running Go tests..."
    go test -v ./...
    
    # Run Python tests if they exist
    if [ -d "tests" ]; then
        print_status "Running Python tests..."
        source venv/bin/activate
        python -m pytest tests/ -v
        deactivate
    fi
}

# Build the application
build_app() {
    print_status "Building the application..."
    
    # Build for current platform
    go build -o bin/job-scraper main_clean.go
    
    # Build for multiple platforms
    print_status "Building for multiple platforms..."
    
    # Linux AMD64
    GOOS=linux GOARCH=amd64 go build -o bin/job-scraper-linux-amd64 main_clean.go
    
    # Linux ARM64
    GOOS=linux GOARCH=arm64 go build -o bin/job-scraper-linux-arm64 main_clean.go
    
    # macOS AMD64
    GOOS=darwin GOARCH=amd64 go build -o bin/job-scraper-darwin-amd64 main_clean.go
    
    # macOS ARM64 (Apple Silicon)
    GOOS=darwin GOARCH=arm64 go build -o bin/job-scraper-darwin-arm64 main_clean.go
    
    # Windows AMD64
    GOOS=windows GOARCH=amd64 go build -o bin/job-scraper-windows-amd64.exe main_clean.go
}

# Create distribution package
create_dist() {
    print_status "Creating distribution package..."
    
    DIST_DIR="dist/job-scraper-$(date +%Y%m%d-%H%M%S)"
    mkdir -p "$DIST_DIR"
    
    # Copy binaries
    cp bin/* "$DIST_DIR/"
    
    # Copy configuration files
    cp -r config "$DIST_DIR/"
    
    # Copy scripts
    cp -r scripts "$DIST_DIR/"
    
    # Copy README
    cp README.md "$DIST_DIR/"
    
    # Create run script
    cat > "$DIST_DIR/run.sh" << EOF
#!/bin/bash
# Job Scraper Run Script

# Get the directory of this script
DIR="\$( cd "\$( dirname "\${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

# Run the appropriate binary
case "\$(uname -s)" in
    Linux*)
        case "\$(uname -m)" in
            x86_64) exec "\$DIR/job-scraper-linux-amd64" ;;
            aarch64) exec "\$DIR/job-scraper-linux-arm64" ;;
            *) echo "Unsupported architecture: \$(uname -m)" ;;
        esac
        ;;
    Darwin*)
        case "\$(uname -m)" in
            x86_64) exec "\$DIR/job-scraper-darwin-amd64" ;;
            arm64) exec "\$DIR/job-scraper-darwin-arm64" ;;
            *) echo "Unsupported architecture: \$(uname -m)" ;;
        esac
        ;;
    CYGWIN*|MINGW*|MSYS*)
        exec "\$DIR/job-scraper-windows-amd64.exe"
        ;;
    *)
        echo "Unsupported OS: \$(uname -s)"
        ;;
esac
EOF
    
    chmod +x "$DIST_DIR/run.sh"
    
    # Create archive
    cd dist
    tar -czf "job-scraper-$(date +%Y%m%d-%H%M%S).tar.gz" "$(basename $DIST_DIR)"
    cd ..
    
    print_status "Distribution package created: dist/job-scraper-$(date +%Y%m%d-%H%M%S).tar.gz"
}

# Clean build artifacts
clean() {
    print_status "Cleaning build artifacts..."
    rm -rf bin/
    rm -rf dist/
    go clean -cache
}

# Main function
main() {
    echo "Job Scraper TUI Application Build Script"
    echo "========================================="
    
    # Parse command line arguments
    case "${1:-build}" in
        "deps")
            check_go
            check_python
            install_go_deps
            install_python_deps
            create_dirs
            create_configs
            ;;
        "test")
            run_tests
            ;;
        "build")
            check_go
            check_python
            install_go_deps
            install_python_deps
            create_dirs
            create_configs
            run_tests
            build_app
            ;;
        "dist")
            check_go
            check_python
            install_go_deps
            install_python_deps
            create_dirs
            create_configs
            run_tests
            build_app
            create_dist
            ;;
        "clean")
            clean
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  deps    - Install dependencies only"
            echo "  test    - Run tests only"
            echo "  build   - Build the application (default)"
            echo "  dist    - Create distribution package"
            echo "  clean   - Clean build artifacts"
            echo "  help    - Show this help message"
            ;;
        *)
            print_error "Unknown command: $1"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
    
    print_status "Done! 🎉"
}

# Run main function with all arguments
main "$@"