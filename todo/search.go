package todo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	dbconnection "todoapp/db"

	"github.com/gorilla/mux"
)

func SearchTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tods []Todos
	var todo Todos

	db = dbconnection.DBConn()

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
