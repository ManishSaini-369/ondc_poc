
package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	RDB *redis.Client
)

func ConnectDB() error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	return nil
}

func ConnectRedis() error {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s", os.Getenv("REDIS_ADDR")))
	if err != nil {
		return err
	}

	RDB = redis.NewClient(opt)
	return nil
}
