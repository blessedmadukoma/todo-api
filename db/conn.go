package db

import (
	"database/sql"
	"fmt"
	"os"
)

func DBConn() (DB *sql.DB) {
	dbDriver := os.Getenv("DIALECT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	DB, err := sql.Open(dbDriver, dbURI)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("DB Connected!!")

	return DB
}
