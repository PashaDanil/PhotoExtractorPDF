#!/bin/bash

# Script to generate OpenAPI documentation and update frontend API client

set -e

# Get the root directory (two levels up from this script)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
JOBS_API_DIR="$ROOT_DIR/apps/jobs-api"
CONTRACTS_DIR="$ROOT_DIR/libs/contracts/openapi"
FRONTEND_DIR="$ROOT_DIR/apps/frontend"

echo "Generating Swagger documentation..."

# Navigate to jobs-api directory and run swag init
cd "$JOBS_API_DIR"
swag init -g main.go -d ./cmd/go-api,./internal/http/handler -o ./docs --outputTypes yaml,json
echo "✓ Swagger documentation generated"

# Copy the generated YAML to contracts directory
echo "Copying OpenAPI spec to contracts..."
SOURCE_YAML="$JOBS_API_DIR/docs/swagger.yaml"
TARGET_YAML="$CONTRACTS_DIR/imgpdf.yaml"

if [ -f "$SOURCE_YAML" ]; then
    cp "$SOURCE_YAML" "$TARGET_YAML"
    echo "✓ OpenAPI spec copied to $TARGET_YAML"
else
    echo "Error: Source YAML file not found: $SOURCE_YAML"
    exit 1
fi

# Generate frontend API client
echo "Generating frontend API client..."
cd "$FRONTEND_DIR"
npm run api:gen
echo "✓ Frontend API client generated"

echo ""
echo "✓ All done!"
