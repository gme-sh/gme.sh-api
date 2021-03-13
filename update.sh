#!/bin/bash

# docker-compose down
echo "docker-compose down"
docker-compose down

echo "git stash, pull, git stash pop"
git stash
git fetch --all
git pull origin main
git stash pop

echo "docker-compose build"
docker-compose build

echo "Done!"