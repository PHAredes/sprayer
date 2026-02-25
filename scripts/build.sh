#!/bin/bash
cd /vercel/share/v0-project
echo "=== Building sprayer ==="
go build ./...
echo "Build exit code: $?"
echo ""
echo "=== Running tests ==="
go test ./src/ui/tui/... -v -count=1 2>&1 | head -200
echo "Test exit code: ${PIPESTATUS[0]}"
