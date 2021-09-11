package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dbconnection "todoapp/db"

	"github.com/gorilla/mux"
)

var todos []Todos

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var todo Todos

	db = dbconnection.DBConn()

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

func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db = dbconnection.DBConn()
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
