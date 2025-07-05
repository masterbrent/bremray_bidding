#!/bin/bash

# Run migrations for Bremray Electrical Bidding App

# Get database URL from environment or use default
DB_URL="${DATABASE_URL:-postgres://postgres@localhost/bremray_dev?sslmode=disable}"

echo "Running migrations..."
echo "Database URL: $DB_URL"

# Function to run a SQL file
run_sql_file() {
    local file=$1
    echo "Running migration: $(basename $file)"
    psql "$DB_URL" -f "$file" 2>&1 | grep -E "(ERROR|NOTICE|CREATE|ALTER|INSERT)" || echo "  ✓ Completed"
}

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
MIGRATIONS_DIR="$SCRIPT_DIR/backend/migrations"

# Check if migrations directory exists
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Error: Migrations directory not found at $MIGRATIONS_DIR"
    exit 1
fi

# Run all SQL files in order
for migration in $(ls -1 "$MIGRATIONS_DIR"/*.sql | sort); do
    run_sql_file "$migration"
done

echo ""
echo "Migrations completed!"
echo ""
echo "Checking template_phases table..."
psql "$DB_URL" -c "SELECT COUNT(*) as phase_count FROM template_phases;" 2>/dev/null || echo "Note: template_phases table might not exist yet"
