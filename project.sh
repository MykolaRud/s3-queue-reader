#!/usr/bin/env bash

# Define Project Name
PROJECT_NAME="queue-reader"

# Connect to Docker Daemon
#eval $(docker-machine env default)

cd .docker

# Execute commands sequence
case "$1" in
	"build")
		# Build Composition Images
		docker-compose -f docker-compose.yml -p ${PROJECT_NAME} build
	;;
	"run")
		# Run Composition
		docker-compose -p ${PROJECT_NAME} up -d
	;;
	"stop")
		# Stop Composition
		docker-compose -p ${PROJECT_NAME} stop
	;;
	"down")
		# Stop Composition and remove its containers
		docker-compose -p ${PROJECT_NAME} down
	;;
	"down-remove")
		# Stop Composition and remove its containers and all volumes
		docker-compose -p ${PROJECT_NAME} down -v --remove-orphans
	;;
	*)
		echo "Invalid argument: $1"
		echo "Available commands: build, run, stop, down"
	;;
esac
