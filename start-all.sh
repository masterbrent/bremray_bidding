#!/bin/bash

# Start all services for the Bremray Electrical Bidding App

echo "Starting Bremray Electrical Bidding App..."

# Check if PostgreSQL is running
if ! pgrep -x "postgres" > /dev/null; then
    echo "PostgreSQL is not running. Please start it first."
    echo "On macOS with Homebrew: brew services start postgresql@14"
    exit 1
fi

# Start backend server
echo "Starting backend server..."
cd backend
# Use DATABASE_URL from environment or .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi
# Default to local dev database if not set
export DATABASE_URL="${DATABASE_URL:-postgres://postgres@localhost/bremray_dev?sslmode=disable}"
./bin/server > server.log 2>&1 &
BACKEND_PID=$!
echo "Backend server started with PID: $BACKEND_PID"

# Wait a bit for backend to start
sleep 2

# Check if backend started successfully
if ! lsof -i :8080 > /dev/null; then
    echo "Backend server failed to start. Check backend/server.log for errors."
    exit 1
fi

# Start frontend server
echo "Starting frontend server..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!
echo "Frontend server started with PID: $FRONTEND_PID"

# Wait for frontend to start
sleep 3

echo ""
echo "✅ All services started successfully!"
echo ""
echo "Access the app at: http://localhost:5173"
echo "Backend API at: http://localhost:8080"
echo ""
echo "To stop all services, run: ./stop-all.sh"
echo ""
echo "PIDs saved to .pids file"

# Save PIDs for stop script
echo "BACKEND_PID=$BACKEND_PID" > ../.pids
echo "FRONTEND_PID=$FRONTEND_PID" >> ../.pids
