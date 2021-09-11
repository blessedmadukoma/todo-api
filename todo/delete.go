package todo

import (
	"encoding/json"
	"log"
	"net/http"
	dbconnection "todoapp/db"

	"github.com/gorilla/mux"
)

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbconnection.DBConn()
	defer db.Close()

	params := mux.Vars(r)

	statement := `DELETE FROM todos WHERE id=$1`

	_, err := db.Exec(statement, params["id"])
	if err != nil {
		log.Fatal("Error deleting record:", err)
	}

	json.NewEncoder(w).Encode("row deleted")
}
