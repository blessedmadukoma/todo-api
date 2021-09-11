package todo

import (
	"encoding/json"
	"fmt"
	"net/http"
	dbconnection "todoapp/db"
)

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbconnection.DBConn()

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
