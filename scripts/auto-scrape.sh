#!/bin/bash

# Auto-scrape script for job-scraper
# This script can be run via cron for automatic job updates

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
LOG_FILE="$PROJECT_ROOT/logs/scrape.log"
DATA_DIR="$HOME/.jobscraper"

# Create directories
mkdir -p "$(dirname "$LOG_FILE")"
mkdir -p "$DATA_DIR"

# Function to log messages
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

# Function to send notification
notify() {
    local message="$1"
    log "NOTIFICATION: $message"
    
    # You can add notification methods here:
    # - Email
    # - Push notification
    # - Desktop notification
    
    # Example: Desktop notification (if available)
    if command -v notify-send >/dev/null 2>&1; then
        notify-send "Job Scraper" "$message"
    fi
}

# Main scraping function
scrape_jobs() {
    log "Starting automated job scraping..."
    
    cd "$PROJECT_ROOT"
    
    # Run the scraper
    if go build -o sprayer ./cmd/sprayer/main.go 2>>"$LOG_FILE"; then
        log "Build successful"
    else
        log "ERROR: Build failed"
        return 1
    fi
    
    # Run scraping
    if ./sprayer 2>>"$LOG_FILE"; then
        log "Scraping completed successfully"
        
        # Count jobs in database
        job_count=$(sqlite3 "$DATA_DIR/sprayer.db" "SELECT COUNT(*) FROM jobs;" 2>/dev/null || echo "0")
        log "Current job count: $job_count"
        
        # Check if we reached target
        if [ "$job_count" -ge 300 ]; then
            notify "🎉 Target reached! Found $job_count jobs"
        else
            notify "Scraping complete. Found $job_count jobs (target: 300)"
        fi
        
        return 0
    else
        log "ERROR: Scraping failed"
        return 1
    fi
}

# Export jobs function
export_jobs() {
    log "Exporting jobs..."
    cd "$PROJECT_ROOT"
    
    # Export to JSON
    timestamp=$(date '+%Y%m%d_%H%M%S')
    export_file="outputs/jobs_export_${timestamp}.json"
    
    mkdir -p outputs
    
    if go build -o sprayer-cli ./cmd/cli/main.go 2>>"$LOG_FILE"; then
        ./sprayer-cli list --keywords "" > "$export_file" 2>>"$LOG_FILE"
        if [ $? -eq 0 ]; then
            log "Jobs exported to $export_file"
        else
            log "ERROR: Export failed"
        fi
    fi
}

# Cleanup old jobs function
cleanup_old_jobs() {
    log "Cleaning up old jobs..."
    
    # Remove jobs older than 30 days
    cutoff_date=$(date -d "30 days ago" '+%Y-%m-%d')
    deleted_count=$(sqlite3 "$DATA_DIR/sprayer.db" "DELETE FROM jobs WHERE posted_date < '$cutoff_date'; SELECT changes();" 2>/dev/null || echo "0")
    
    log "Removed $deleted_count old jobs"
}

# Main execution
main() {
    case "${1:-scrape}" in
        "scrape")
            scrape_jobs
            ;;
        "export")
            export_jobs
            ;;
        "cleanup")
            cleanup_old_jobs
            ;;
        "all")
            scrape_jobs
            export_jobs
            cleanup_old_jobs
            ;;
        *)
            echo "Usage: $0 {scrape|export|cleanup|all}"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"