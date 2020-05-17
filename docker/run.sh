#!/usr/bin/env bash

# Check if docker is running
docker_state=$(docker info >/dev/null 2>&1)
if [[ $? -ne 0 ]]; then
    echo "Docker not running. "
    exit 1
fi

# build and
docker-compose rm -f && docker-compose build && docker-compose up
