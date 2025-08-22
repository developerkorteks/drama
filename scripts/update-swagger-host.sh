#!/bin/bash

# Script to update Swagger host dynamically
# Usage: ./scripts/update-swagger-host.sh [domain]

DOMAIN=${1:-localhost:52983}
DOCS_FILE="docs/docs.go"

if [ ! -f "$DOCS_FILE" ]; then
    echo "Error: $DOCS_FILE not found. Please run 'swag init' first."
    exit 1
fi

echo "Updating Swagger host to: $DOMAIN"

# Update the host in docs.go
sed -i.bak "s/\"host\": \"[^\"]*\"/\"host\": \"$DOMAIN\"/" "$DOCS_FILE"

if [ $? -eq 0 ]; then
    echo "Successfully updated Swagger host to $DOMAIN"
    rm -f "$DOCS_FILE.bak"
else
    echo "Error updating Swagger host"
    exit 1
fi