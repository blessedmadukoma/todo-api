package main

import (
	"go-gorm-pg/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error connecting to main .env:", err)
	}

	// DB := models.DBCon()

	// DB.Create(&models.User{
	// 	Firstname: "Jimmy", Lastname: "Johnson", Email: "jimmyjohnson@gmail.com", Password: "fri123456", Todos: []models.Todo{
	// 		{Text: "First 1 Todo"},
	// 		{Text: "Second 1 Todo"},
	// 		{Text: "Third 1 Todo"},
	// 	},
	// })
	// DB.Create(&models.User{
	// 	Firstname: "Howard", Lastname: "Hills", Email: "howardhills@gmail.com", Password: "fri123456", Todos: []models.Todo{
	// 		{Text: "First 2 Todo"},
	// 		{Text: "Second 2 Todo"},
	// 		{Text: "Third 2 Todo"},
	// 	},
	// })
	// DB.Create(&models.User{
	// 	Firstname: "Craig", Lastname: "Colbin", Email: "craigcolbin@gmail.com", Password: "fri123456", Todos: []models.Todo{
	// 		{Text: "First 3 Todo"},
	// 		{Text: "Second 3 Todo"},
	// 		{Text: "Third 3 Todo"},
	// 	},
	// })

	port := os.Getenv("PORT")
	log.Println("Server starting on port:", port)

	router := routes.Router()
	log.Fatal(http.ListenAndServe(":"+port, router))

}
