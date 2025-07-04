#!/bin/bash

# Stop all services for the Bremray Electrical Bidding App

echo "Stopping Bremray Electrical Bidding App..."

# Load PIDs if file exists
if [ -f .pids ]; then
    source .pids
    
    # Stop backend
    if [ ! -z "$BACKEND_PID" ] && kill -0 $BACKEND_PID 2>/dev/null; then
        echo "Stopping backend server (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
    fi
    
    # Stop frontend
    if [ ! -z "$FRONTEND_PID" ] && kill -0 $FRONTEND_PID 2>/dev/null; then
        echo "Stopping frontend server (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
    fi
    
    rm .pids
else
    echo "No .pids file found. Searching for processes..."
    
    # Find and kill backend server
    BACKEND_PID=$(lsof -t -i :8080)
    if [ ! -z "$BACKEND_PID" ]; then
        echo "Stopping backend server (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
    fi
    
    # Find and kill frontend server  
    FRONTEND_PID=$(lsof -t -i :5173)
    if [ ! -z "$FRONTEND_PID" ]; then
        echo "Stopping frontend server (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
    fi
fi

echo "✅ All services stopped!"
