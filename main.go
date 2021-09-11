package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Todos struct {
	ID   string `json:"id"`
	Item string `json:"item"`
}

var todos []Todos

var (
	db  *sql.DB
	err error
)

func dbConn() (db *sql.DB) {
	dbDriver := os.Getenv("DIALECT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open(dbDriver, dbURI)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("DB Connected!!")

	return db
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todos

	db = dbConn()

	rows, err := db.Query("SELECT id, item FROM todos")
	if err != nil {
		fmt.Println("Error when inserting: ", err.Error())
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&todo.ID, &todo.Item)
		if err != nil {
			log.Println("Error occured during getting todos:", err)
		}
		todos = append(todos, todo)
	}
	json.NewEncoder(w).Encode(todos)
	defer db.Close()
	defer rows.Close()
}

func getTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db = dbConn()
	statement := `SELECT id, item FROM todos WHERE id=$1`

	var todo Todos
	params := mux.Vars(r)

	defer db.Close()

	rows := db.QueryRowContext(context.TODO(), statement, params["id"])
	err = rows.Scan(&todo.ID, &todo.Item)
	if err != nil {
		log.Println("Error getting single ID:", err)
	}
	// fmt.Println(*todo)
	json.NewEncoder(w).Encode(todo)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db = dbConn()

	var todo Todos
	json.NewDecoder(r.Body).Decode(&todo)

	defer db.Close()

	statement := `INSERT INTO todos(item) VALUES($1) RETURNING id, item`
	err := db.QueryRow(statement, todo.Item).Scan(&todo.ID, &todo.Item)
	if err != nil {
		fmt.Println("Error when inserting: ", err)
		panic(err.Error())
	}
	fmt.Println(todo)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var todo Todos
	json.NewDecoder(r.Body).Decode(&todo)

	db := dbConn()
	defer db.Close()

	statement := `UPDATE todos SET item=$2 WHERE id=$1`

	_, err = db.Exec(statement, params["id"], todo.Item)
	if err != nil {
		log.Fatal("Error updating record:", err)
	}

	statement = `SELECT id, item FROM todos WHERE id=$1`

	rows := db.QueryRowContext(context.TODO(), statement, params["id"])
	err = rows.Scan(&todo.ID, &todo.Item)
	if err != nil {
		log.Fatal("Error getting new record:", err)
	}

	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbConn()
	defer db.Close()

	params := mux.Vars(r)

	statement := `DELETE FROM todos WHERE id=$1`

	_, err = db.Exec(statement, params["id"])
	if err != nil {
		log.Fatal("Error deleting record:", err)
	}

	json.NewEncoder(w).Encode("row deleted")
}

func searchTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tods []Todos
	var todo Todos

	db = dbConn()

	params := mux.Vars(r)
	fmt.Println(params["key"])
	statement := `SELECT id, item FROM todos WHERE item LIKE '%' || $1 || '%'`

	rows, err := db.Query(statement, params["key"])

	if err != nil {
		fmt.Println("Error when searching: ", err)
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&todo.ID, &todo.Item)
		if err != nil {
			log.Println("Error occured getting todos:", err)
		}
		fmt.Println(todo)
		tods = append(tods, todo)
	}
	json.NewEncoder(w).Encode(tods)
	defer db.Close()
	defer rows.Close()
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()

	router.HandleFunc("/todos", getTodos).Methods("GET")
	router.HandleFunc("/todo/create", createTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", getTodo).Methods("GET")
	router.HandleFunc("/todo/{id}", updateTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", deleteTodo).Methods("DELETE")
	router.HandleFunc("/todo/search/{key}", searchTodos).Methods("POST")

	port := os.Getenv("PORT")
	fmt.Println("Server starting at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
