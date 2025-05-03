#!/usr/bin/env bash

# Exit immediately if a command exits with a non-zero status
set -e

# Capture errors in pipes too
set -o pipefail

env_file=${1:-./.env.local.dev}
compose_file=${2:-./docker-compose.dev.yml}

# Check if files exist
if [ ! -f "$env_file" ]; then
    echo "Error: Environment file '$env_file' not found" >&2
    exit 1
fi

if [ ! -f "$compose_file" ]; then
    echo "Error: Docker compose file '$compose_file' not found" >&2
    exit 1
fi

# echo "compose file: $compose_file"
export DEV=1

# Function to handle errors
handle_error() {
    echo "Error occurred at line $1, exiting..." >&2
    stop_docker
    exit 1
}

# Set up trap for error handling
trap 'handle_error $LINENO' ERR

cleanup() {
    echo "Ctrl+C pressed. Cleaning up..."
    stop_docker
    exit 1
}

# Source required scripts with error checking

# source ./scripts/init_env.sh $env_file
source ./scripts/encrypt_user_pass.sh "$env_file" || { echo "Failed to execute encrypt_user_pass.sh" >&2; exit 1; }
source ./scripts/run_docker.sh "$env_file" "$compose_file" || { echo "Failed to execute run_docker.sh" >&2; exit 1; }
source ./scripts/pre_test.sh || { echo "Failed to execute pre_test.sh" >&2; exit 1; }

run_docker || { echo "Failed to run docker" >&2; exit 1; }
trap cleanup SIGINT
fresh || { echo "Failed to run fresh command" >&2; exit 1; }
