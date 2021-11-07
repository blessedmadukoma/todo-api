package routes

import (
	"go-gorm-pg/handlers"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Home).Methods("GET")
	router.HandleFunc("/signup", handlers.Signup).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")

	router.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	router.HandleFunc("/todo/create", handlers.CreateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", handlers.GetTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", handlers.EditTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", handlers.DeleteTodo).Methods("DELETE")
	router.HandleFunc("/todo/search/{key}", handlers.SearchTodos).Methods("POST")
	return router
}
