#!/bin/bash

# Kill any existing process on port 8080
echo "Checking for existing processes on port 8080..."
lsof -ti:8080 | xargs kill -9 2>/dev/null || true

# Set environment variables if needed
# Load from .env file if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi
# Use defaults if not set
export DATABASE_URL="${DATABASE_URL:-postgres://postgres@localhost/bremray_dev?sslmode=disable}"
export PORT="${PORT:-8080}"

# Start the server
echo "Starting Bremray backend server on port 8080..."
cd "$(dirname "$0")"
go run cmd/server/main.go