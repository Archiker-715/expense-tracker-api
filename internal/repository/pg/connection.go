package pg

import (
	"fmt"
	"log"
	"os"

	"github.com/Archiker-715/expense-tracker-api/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := DB.AutoMigrate(&entity.Expense{}, &entity.Users{}); err != nil {
		log.Fatalf("failed to migrare DB: %v", err)
	}

	log.Println("Database connected successfully")
}
