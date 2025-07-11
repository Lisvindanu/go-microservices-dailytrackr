package database

import (
	"database/sql"
	"fmt"
	"log"

	"dailytrackr/shared/config"
	_ "github.com/go-sql-driver/mysql"
)

// GetMySQLConnection creates and returns a MySQL database connection
func GetMySQLConnection(cfg *config.Config) (*sql.DB, error) {
	// Use the new GetMySQLDSN method
	dsn := cfg.GetMySQLDSN()

	log.Printf("Connecting to MySQL: %s@%s:%s/%s",
		cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	log.Println("✅ Successfully connected to MySQL database")
	return db, nil
}

// TestMySQLConnection tests MySQL database connectivity
func TestMySQLConnection() {
	cfg := config.LoadConfig()

	db, err := GetMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	defer db.Close()

	// Test query
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}

	fmt.Printf("✅ MySQL Version: %s\n", version)

	// Test table existence
	var tableCount int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?", cfg.DBName).Scan(&tableCount)
	if err != nil {
		log.Printf("⚠️  Warning: Could not check tables: %v", err)
	} else {
		fmt.Printf("📊 Number of tables in database: %d\n", tableCount)
	}
}
