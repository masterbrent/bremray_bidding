#!/bin/bash

# Test adding phases to a template

API_URL="http://localhost:8080/api"

# First, get all templates to find one to test with
echo "Getting templates..."
curl -s "$API_URL/templates" | jq -r '.[] | "\(.id) - \(.name)"'

echo ""
echo "Testing phase update on first template..."

# Get the first template ID
TEMPLATE_ID=$(curl -s "$API_URL/templates" | jq -r '.[0].id')
TEMPLATE_NAME=$(curl -s "$API_URL/templates" | jq -r '.[0].name')

if [ -z "$TEMPLATE_ID" ]; then
    echo "No templates found!"
    exit 1
fi

echo "Using template: $TEMPLATE_ID - $TEMPLATE_NAME"
echo ""

# Update template with phases
echo "Adding phases to template..."
curl -X PUT "$API_URL/templates/$TEMPLATE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "'"$TEMPLATE_NAME"'",
    "description": "",
    "isActive": true,
    "phases": [
      {"name": "Rough", "order": 1, "description": "Rough electrical work"},
      {"name": "Final", "order": 2, "description": "Final connections and testing"}
    ]
  }' -s | jq '.'

echo ""
echo "Verifying phases were saved..."
curl -s "$API_URL/templates/$TEMPLATE_ID" | jq '.phases'
