#!/bin/bash

# Stop development services

echo "🛑 Stopping ITMS Development Services..."

if [ -f .pids ]; then
    PIDS=$(cat .pids)
    for PID in $PIDS; do
        if ps -p $PID > /dev/null 2>&1; then
            kill $PID
            echo "  Stopped process $PID"
        fi
    done
    rm .pids
else
    # Fallback: kill by process name
    pkill -f "go run.*auth-service"
    pkill -f "go run.*tenant-service"
    pkill -f "go run.*ticket-service"
    pkill -f "go run.*asset-service"
    pkill -f "go run.*gateway"
fi

echo "✅ All services stopped"