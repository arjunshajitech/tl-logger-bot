#!/bin/bash

# Stop and remove containers, networks, volumes, and images created by docker-compose up
echo "Stopping and removing containers..."
sudo docker compose down