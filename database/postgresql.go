package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var PG *gorm.DB

func init() {
	host := getEnv("DB_HOST", "101.201.49.155")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "secretpassword")
	dbname := getEnv("DB_NAME", "mydatabase")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		host, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	PG = db
}
