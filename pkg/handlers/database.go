package handlers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the PostgreSQL database connection
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Create table if it doesn't exist
	if err = createTouristTable(); err != nil {
		return fmt.Errorf("failed to create tourist table: %v", err)
	}

	return nil
}

// createTouristTable creates the tourist_data table
func createTouristTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS tourist_data (
		id VARCHAR(36) PRIMARY KEY,
		tourist_name VARCHAR(255) NOT NULL,
		digital_expiry VARCHAR(50) NOT NULL,
		data_hash VARCHAR(64) NOT NULL,
		raw_data JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_tourist_data_id ON tourist_data(id);
	CREATE INDEX IF NOT EXISTS idx_tourist_data_hash ON tourist_data(data_hash);`

	_, err := db.Exec(query)
	return err
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
