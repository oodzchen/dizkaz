#!/usr/bin/env bash

# Error handling - script will exit immediately if a command fails
set -e

# Function for logging information
log_info() {
    echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Function for logging warnings
log_warn() {
    echo "[WARN] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Function for logging errors
log_error() {
    echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') - $1"
}

# Error handler function
handle_error() {
    log_error "Error occurred at line $1, exit status: $2"
    exit $2
}

# Set error trap
trap 'handle_error ${LINENO} $?' ERR

# Check if command was successful
check_command() {
    if [ $? -ne 0 ]; then
        log_error "$1 failed"
        return 1
    else
        log_info "$1 successful"
        return 0
    fi
}

# Initialize variables
env_file_name=${1:-.env.local.dev}
GIT_REPO_PATH="$HOME/go/src/github.com/oodzchen/dizkaz/"
DOCKER_COMPOSE_PATH="$GIT_REPO_PATH/docker-compose.yml"
ENV_FILE="$GIT_REPO_PATH/$env_file_name"

# Check if necessary files exist
log_info "Checking necessary files..."
if [ ! -f "$ENV_FILE" ]; then
    log_error "Environment file $ENV_FILE does not exist"
    exit 1
fi

if [ ! -f "$DOCKER_COMPOSE_PATH" ]; then
    log_error "Docker Compose file $DOCKER_COMPOSE_PATH does not exist"
    exit 1
fi

# Set application version
if [ -z "$APP_VERSION" ]; then
    APP_VERSION="latest"
    log_warn "APP_VERSION not specified, using default: $APP_VERSION"
else
    log_info "Using specified APP_VERSION: $APP_VERSION"
fi

# Change to project directory
log_info "Changing to project directory: $GIT_REPO_PATH"
cd $GIT_REPO_PATH || { log_error "Cannot change to project directory"; exit 1; }

# Pull latest code
log_info "Pulling latest code from Git..."
git pull
check_command "Git code pull"

# Update environment file with app version
log_info "Updating app version in environment file..."
sed -i "s/APP_VERSION=.*/APP_VERSION=$APP_VERSION/" $ENV_FILE
check_command "Environment file update"

# Pull latest Docker images
log_info "Pulling latest Docker images..."
docker compose --env-file $ENV_FILE -f $DOCKER_COMPOSE_PATH pull
check_command "Docker image pull"

# Show currently running services
log_info "Currently running services:"
docker compose --env-file $ENV_FILE -f $DOCKER_COMPOSE_PATH ps

# Restart all Docker Compose services
log_info "Restarting all Docker Compose services..."
docker compose --env-file $ENV_FILE -f $DOCKER_COMPOSE_PATH up -d
check_command "Service startup"

# Verify services are running properly
log_info "Verifying service status..."
docker compose --env-file $ENV_FILE -f $DOCKER_COMPOSE_PATH ps
running_services=$(docker compose --env-file $ENV_FILE -f $DOCKER_COMPOSE_PATH ps --services --filter "status=running")
if [ -z "$running_services" ]; then
    log_error "No services are running, deployment may have failed"
    docker compose --env-file $ENV_FILE -f $DOCKER_COMPOSE_PATH logs
    exit 1
fi

# Clean up unused images
log_info "Cleaning unused images..."
docker image prune -af
check_command "Unused image cleanup"

log_info "Deployment completed successfully"
log_info "Running services: $running_services"
