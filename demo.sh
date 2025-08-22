#!/bin/bash

echo "=== DramaQu API Demo ==="
echo ""

# Start the server in background
echo "Starting DramaQu API server..."
go run main.go &
SERVER_PID=$!

# Wait for server to start
echo "Waiting for server to start..."
sleep 5

echo ""
echo "=== Testing Health Endpoint ==="
curl -s http://localhost:8080/health | jq '.'

echo ""
echo "=== Testing Home Endpoint ==="
echo "Note: This may take a few seconds as it scrapes data from dramaqu.ad..."
curl -s http://localhost:8080/api/v1/home | jq '.'

echo ""
echo "=== Swagger Documentation ==="
echo "Visit: http://localhost:8080/swagger/index.html"

echo ""
echo "Press any key to stop the server..."
read -n 1 -s

# Kill the server
kill $SERVER_PID
echo "Server stopped."