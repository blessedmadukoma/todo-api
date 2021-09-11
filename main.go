package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	todoo "todoapp/todo"

	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/todos", todoo.GetTodos).Methods("GET")
	router.HandleFunc("/todo/create", todoo.CreateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", todoo.GetTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", todoo.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", todoo.DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todo/search/{key}", todoo.SearchTodos).Methods("POST")

	port := os.Getenv("PORT")
	fmt.Println("Server starting at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
