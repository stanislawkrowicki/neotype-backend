package mysql

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"neotype-backend/pkg/config"
	"os"
)

const (
	envDatabase = "MYSQL_DATABASE"
	envUser     = "MYSQL_USER"
	envPassword = "MYSQL_PASSWORD"
	configAddr  = "addr"
	configPort  = "port"
	charset     = "utf8mb4"
	parseTime   = true
	locale      = "Local"
)

func generateDsn() string {
	_ = godotenv.Load()

	addr, err := config.Get("mysql", configAddr)
	port, err := config.Get("mysql", configPort)
	database := os.Getenv(envDatabase)
	user := os.Getenv(envUser)
	password := os.Getenv(envPassword)

	if err != nil {
		panic("Failed to load config for mysql.")
	}

	dsn := fmt.Sprintf("%s:%s", user, password) +
		fmt.Sprintf("@tcp(%s:%s)", addr, port) +
		fmt.Sprintf("/%s?", database) +
		fmt.Sprintf("charset=%s", charset) +
		fmt.Sprintf("&parseTime=%t", parseTime) +
		fmt.Sprintf("&loc=%s", locale)

	return dsn
}

func NewConnection() *gorm.DB {
	dsn := generateDsn()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to open connection to database. %v", err))
	}

	return db
}
