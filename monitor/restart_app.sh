#!/bin/sh
# restart_app.sh

SERVICE_NAME="web-app"

# Проверяем, существует ли Docker CLI
if ! command -v docker > /dev/null; then
    echo "FATAL: Docker client not found inside container."
    exit 1
fi

if [ ! -S /var/run/docker.sock ]; then
    echo "FATAL: Docker socket /var/run/docker.sock not mounted."
    exit 1
fi

# Ищем ID контейнера по метке (label) 'com.docker.compose.service'
CONTAINER_ID=$(docker ps -q -f "label=com.docker.compose.service=$SERVICE_NAME" -f "status=running")

if [ -z "$CONTAINER_ID" ]; then
    echo "WARNING: Service $SERVICE_NAME container is not running (ID not found). Assuming it crashed or stopped."
    
    # Пробуем запустить, если он просто остановлен
    STOPPED_CONTAINER_ID=$(docker ps -aq -f "label=com.docker.compose.service=$SERVICE_NAME" -f "status=exited")
    if [ ! -z "$STOPPED_CONTAINER_ID" ]; then
        echo "INFO: Found stopped container $STOPPED_CONTAINER_ID. Starting it."
        docker start "$STOPPED_CONTAINER_ID"
    fi
    
    exit 0
else
    # Перезапускаем, если найден ID запущенного контейнера
    echo "INFO: Restarting running container $CONTAINER_ID for service $SERVICE_NAME."
    docker restart "$CONTAINER_ID"
    exit 0
fi