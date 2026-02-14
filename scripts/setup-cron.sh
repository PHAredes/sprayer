#!/bin/bash

# Cron job setup script for job-scraper
# This script helps set up automated scraping via cron

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CRON_SCRIPT="$SCRIPT_DIR/auto-scrape.sh"

# Check if script exists
if [ ! -f "$CRON_SCRIPT" ]; then
    echo "Error: Auto-scrape script not found at $CRON_SCRIPT"
    exit 1
fi

# Function to add cron job
add_cron_job() {
    local schedule="$1"
    local command="$2"
    
    # Check if cron job already exists
    if crontab -l 2>/dev/null | grep -q "$command"; then
        echo "Cron job already exists:"
        crontab -l | grep "$command"
    else
        # Add the cron job
        (crontab -l 2>/dev/null; echo "$schedule $command") | crontab -
        echo "Cron job added: $schedule $command"
    fi
}

# Function to show current cron jobs
show_cron_jobs() {
    echo "Current cron jobs for job-scraper:"
    crontab -l 2>/dev/null | grep "auto-scrape.sh" || echo "No job-scraper cron jobs found"
}

# Function to remove all job-scraper cron jobs
remove_cron_jobs() {
    if crontab -l 2>/dev/null | grep -q "auto-scrape.sh"; then
        crontab -l | grep -v "auto-scrape.sh" | crontab -
        echo "All job-scraper cron jobs removed"
    else
        echo "No job-scraper cron jobs found"
    fi
}

# Main function
main() {
    case "${1:-help}" in
        "setup")
            echo "Setting up automated job scraping..."
            
            # Add daily scraping at 6 AM
            add_cron_job "0 6 * * *" "$CRON_SCRIPT scrape"
            
            # Add weekly export on Sunday at 5 AM
            add_cron_job "0 5 * * 0" "$CRON_SCRIPT export"
            
            # Add monthly cleanup on 1st of month at 4 AM
            add_cron_job "0 4 1 * *" "$CRON_SCRIPT cleanup"
            
            echo ""
            echo "Cron jobs setup complete!"
            echo "- Daily scraping: 6 AM"
            echo "- Weekly export: Sunday 5 AM" 
            echo "- Monthly cleanup: 1st of month 4 AM"
            ;;
        "show")
            show_cron_jobs
            ;;
        "remove")
            remove_cron_jobs
            ;;
        "test")
            echo "Running test scrape..."
            "$CRON_SCRIPT" scrape
            ;;
        "help"|"*")
            echo "Usage: $0 {setup|show|remove|test}"
            echo ""
            echo "Commands:"
            echo "  setup   - Set up automated scraping cron jobs"
            echo "  show    - Show current cron jobs"
            echo "  remove  - Remove all job-scraper cron jobs"
            echo "  test    - Run a test scrape"
            echo ""
            echo "Example schedules:"
            echo "  Daily at 6 AM: 0 6 * * *"
            echo "  Every 4 hours: 0 */4 * * *"
            echo "  Weekdays at 8 AM: 0 8 * * 1-5"
            ;;
    esac
}

# Check if crontab is available
if ! command -v crontab >/dev/null 2>&1; then
    echo "Error: crontab command not found. Please install cron daemon."
    exit 1
fi

# Run main function
main "$@"