#!/bin/bash

# Stop and remove containers, networks, volumes, and images created by docker-compose up
echo "Stopping and removing containers..."
sudo docker compose down

# Rebuild and restart the containers in detached mode
echo "Rebuilding and starting containers..."
sudo docker compose up --build -d

echo "Containers are up and running!"
