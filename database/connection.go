package database

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Db *gorm.DB
}

// NewDatabase is the constructor for the database
func NewDatabase() (*Database, error) {
	// Get database configuration from environment variables or use defaults
	host := os.Getenv(fmt.Sprintf("DB_HOST_%s", strings.ToUpper(os.Getenv("ENV"))))
	user := os.Getenv(fmt.Sprintf("DB_USER_%s", strings.ToUpper(os.Getenv("ENV"))))
	password := os.Getenv(fmt.Sprintf("DB_PASSWORD_%s", strings.ToUpper(os.Getenv("ENV"))))
	dbname := os.Getenv(fmt.Sprintf("DB_NAME_%s", strings.ToUpper(os.Getenv("ENV"))))
	port := os.Getenv(fmt.Sprintf("DB_PORT_%s", strings.ToUpper(os.Getenv("ENV"))))

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		Db: db,
	}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.Db
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.Db.AutoMigrate(models...)
}
