#!/bin/bash
set -e

echo "Stopping containers..."
docker compose down

echo "Rebuilding and starting containers..."
docker compose up --build -d

echo "Containers restarted successfully, database data preserved."
