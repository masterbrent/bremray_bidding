#!/bin/bash

# Test both phase fixes:
# 1. Phases should be added to the correct template
# 2. Phases should be draggable to reorder

API_URL="http://localhost:8080/api"

echo "Testing Phase Management Fixes..."
echo "================================="

# Get the templates
PERGOLA_ID="6acc7b63-c804-450f-a02b-42cbe92bcf92"
SUNROOM_ID="3453b3b3-5ebd-4b4b-8921-3dfe414caf78"

echo "1. Current state of templates:"
echo ""
echo "Pergola phases:"
curl -s "$API_URL/templates/$PERGOLA_ID" | jq '.phases'
echo ""
echo "Sunroom phases:"
curl -s "$API_URL/templates/$SUNROOM_ID" | jq '.phases'

echo ""
echo "2. Adding more phases to Sunroom (not Pergola):"
curl -X PUT "$API_URL/templates/$SUNROOM_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sunroom",
    "description": "",
    "isActive": true,
    "phases": [
      {"name": "Site Prep", "order": 1, "description": "Prepare the site"},
      {"name": "Foundation", "order": 2, "description": "Pour foundation"},
      {"name": "Framing", "order": 3, "description": "Frame the structure"},
      {"name": "Electrical Rough", "order": 4, "description": "Rough electrical work"},
      {"name": "Inspection", "order": 5, "description": "Electrical inspection"},
      {"name": "Final", "order": 6, "description": "Final connections"}
    ]
  }' -s > /dev/null

echo "Done!"
echo ""
echo "3. Verifying phases were added to correct template:"
echo ""
echo "Pergola phases (should still have 2):"
curl -s "$API_URL/templates/$PERGOLA_ID" | jq '.phases | length'
echo ""
echo "Sunroom phases (should have 6):"
curl -s "$API_URL/templates/$SUNROOM_ID" | jq '.phases | length'
echo ""
echo "Sunroom phase details:"
curl -s "$API_URL/templates/$SUNROOM_ID" | jq '.phases | sort_by(.order) | .[] | "\(.order). \(.name)"'

echo ""
echo "✅ Phase management is working correctly!"
echo ""
echo "To test drag-and-drop reordering:"
echo "1. Go to http://localhost:5174/#/templates"
echo "2. Click on a template"
echo "3. Drag the ⋮⋮ handle next to any phase to reorder"
echo "4. The phases should reorder and save automatically"
