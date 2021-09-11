package todo

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	dbconnection "todoapp/db"

	"github.com/gorilla/mux"
)

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var todo Todos
	json.NewDecoder(r.Body).Decode(&todo)

	db := dbconnection.DBConn()
	defer db.Close()

	statement := `UPDATE todos SET item=$2 WHERE id=$1`

	_, err := db.Exec(statement, params["id"], todo.Item)
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
