package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres@localhost/bremray_dev?sslmode=disable"
	}

	fmt.Printf("Connecting to database...\n")
	
	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database")

	// Get migrations directory
	migrationsDir := filepath.Join(".", "migrations")
	
	// Read all SQL files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	fmt.Printf("\nFound %d migration files\n", len(sqlFiles))

	// Run each migration
	for _, filename := range sqlFiles {
		filepath := filepath.Join(migrationsDir, filename)
		fmt.Printf("\nRunning migration: %s\n", filename)
		
		// Read file
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Printf("Failed to read %s: %v", filename, err)
			continue
		}

		// Execute SQL
		_, err = db.Exec(string(content))
		if err != nil {
			// Check if it's a "already exists" error
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("  → Already applied (skipping)\n")
			} else {
				log.Printf("  → Error: %v", err)
			}
		} else {
			fmt.Printf("  → Success\n")
		}
	}

	// Check template_phases table
	fmt.Println("\nChecking template_phases table...")
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM template_phases").Scan(&count)
	if err != nil {
		fmt.Printf("Error checking template_phases: %v\n", err)
	} else {
		fmt.Printf("template_phases table has %d rows\n", count)
	}

	// Check if the specific tables exist
	checkTable := func(tableName string) {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)`, tableName).Scan(&exists)
		
		if err != nil {
			fmt.Printf("Error checking %s: %v\n", tableName, err)
		} else if exists {
			fmt.Printf("✓ Table %s exists\n", tableName)
		} else {
			fmt.Printf("✗ Table %s does NOT exist\n", tableName)
		}
	}

	fmt.Println("\nChecking critical tables:")
	checkTable("template_phases")
	checkTable("jobs")
	checkTable("job_templates")
	
	fmt.Println("\nMigration check complete!")
}
