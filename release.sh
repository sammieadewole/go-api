#!/bin/bash
set -e

docker-compose --env-file .env up -d

# sleep for a few seconds
sleep 3

echo "Starting Go App..."
./main
echo "App running at: http://localhost:8080"
