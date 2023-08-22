#!/usr/bin/env sh

echo "Waiting for MySQL to start..."
CONTAINER=$(docker ps --filter "name=mysql" --format "{{.ID}}")
while ! docker exec -it "${CONTAINER}" mysql --user=sqx --password=sqx -e "SELECT 1"; do
	CONTAINER=$(docker ps --filter "name=mysql" --format "{{.ID}}")
	sleep 1
done
echo "MySQL started"
