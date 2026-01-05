package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to Postgres Database
func ConnectPG() {

	dsn := os.Getenv("POSTGRES_URL")

	if dsn == "" {
		log.Fatal("POSTGRES_URL is not set")
	}

	// Create a postgres session
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Postgres connected succesfully")

	DB = db
}

// Migrates all models to database
//
// - models: Models
func MigratePG(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("POSTGRES is not connected")
	}
	return DB.AutoMigrate(models...)
}
