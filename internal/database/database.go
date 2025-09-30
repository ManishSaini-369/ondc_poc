package database

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ondc-poc/internal/models"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
)

// ConnectDB initializes PostgreSQL using GORM
func ConnectDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Auto migrate your models
	if err := db.AutoMigrate(&models.SearchRequest{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	DB = db
	log.Println("PostgreSQL connected and migrated successfully")
}

// ConnectRedis initializes Redis client
func ConnectRedis() error {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s", os.Getenv("REDIS_ADDR")))
	if err != nil {
		return err
	}

	RDB = redis.NewClient(opt)
	log.Println("Redis connected successfully")
	return nil
}
