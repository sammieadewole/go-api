#!/bin/bash
set -e

# echo "Building Go app..."
# docker build -t my-go-app .

echo "Starting Docker Compose with runtime env..."
docker-compose --env-file .env up -d

echo "Databases are up!"
echo "Postgres: $POSTGRES_URL"
echo "MongoDB: $MONGODB_URI"
echo "Mongo Express running at: http://localhost:8081"

# sleep for a few seconds
sleep 3

echo "Starting Go App..."
./main.exe
echo "App running at: http://localhost:8080"
