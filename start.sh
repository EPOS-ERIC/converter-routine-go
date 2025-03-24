#!/bin/sh


export BROKER_HOST="rabbitmq"
export BROKER_USERNAME="changeme"
export BROKER_PASSWORD="changeme"
export BROKER_VHOST="changeme"
export POSTGRESQL_HOST="metadata-catalogue"
export POSTGRESQL_PORT=5433
export POSTGRES_USER="cerif_admin"
export POSTGRESQL_PASSWORD="brgm"
export POSTGRES_DB="cerif"
export PERSISTENCE_NAME="EPOSDataModel"
export PERSISTENCE_NAME_PROCESSING="EPOSProcessing"
export POSTGRESQL_CONNECTION_STRING="jdbc:postgresql://localhost:${POSTGRESQL_PORT},localhost:5432/${POSTGRES_DB}?user=${POSTGRES_USER}&password=${POSTGRESQL_PASSWORD}"

echo "$POSTGRESQL_CONNECTION_STRING"

go run .
