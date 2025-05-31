package database

import (
	"log"
	"os"
	"path/filepath"
)

// InitDB initializes the database with the schema
func InitDB() error {
	log.Println("Initializing database...")
	
	// Read schema file
	schemaPath := filepath.Join("database", "schema.sql")
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return err
	}
	
	schema := string(schemaBytes)
	
	// Execute schema
	_, err = DB.Exec(schema)
	if err != nil {
		return err
	}
	
	log.Println("Database initialized successfully")
	return nil
} 