#!/bin/bash

# Test Phase Management

API_URL="http://localhost:8080/api"
TEMPLATE_ID="YOUR_TEMPLATE_ID_HERE"

echo "Testing Phase Management..."

# Test 1: Get template
echo -e "\n1. Getting template..."
curl -s "$API_URL/templates/$TEMPLATE_ID" | jq '.'

# Test 2: Update template with phases
echo -e "\n2. Updating template with phases..."
curl -X PUT "$API_URL/templates/$TEMPLATE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sunroom",
    "description": "Sunroom installation",
    "phases": [
      {"name": "Rough", "order": 1, "description": "Rough electrical work"},
      {"name": "Final", "order": 2, "description": "Final connections and testing"}
    ]
  }' | jq '.'

# Test 3: Get template again to verify phases
echo -e "\n3. Getting template again to verify phases..."
curl -s "$API_URL/templates/$TEMPLATE_ID" | jq '.'
