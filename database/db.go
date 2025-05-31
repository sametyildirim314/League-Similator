package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sametyildirim314/insider_case/config"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	
	// Get configuration
	cfg := config.GetConfig()
	
	// Connection string with additional options
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)
	
	// Connect to database
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// Set connection pool parameters
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	
	// Check connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	log.Println("Successfully connected to database")
} 