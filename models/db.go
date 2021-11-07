package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBCon() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error connecting to .env:", err)
	}

	dbDriver := os.Getenv("DB_Driver")
	dbUser := os.Getenv("DB_User")
	dbPass := os.Getenv("DB_Password")
	dbName := os.Getenv("DB_Name")
	dbHost := os.Getenv("DB_Host")
	dbPort := os.Getenv("DB_Port")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	// DB, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	sqlDB, err := sql.Open(dbDriver, dbURI)
	if err != nil {
		log.Fatal("Error connecting to sql postgres .env:", err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	DB.AutoMigrate(&User{}, &Todo{})
	fmt.Println("DB Connected!!")

	return DB
}
