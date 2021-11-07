package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-gorm-pg/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	db = models.DBCon()

	todo  models.Todo
	todos []models.Todo
	err   error
	// user  *models.User
)

func GetTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// db.Raw("SELECT * FROM todo").Scan(&todos)

	// json.NewEncoder(w).Encode(todos)

	fmt.Println("Response:", Response.Data["user_id"])
	var count int64 = 0
	m := make(map[string]interface{})
	m["user_id"] = Response.Data["user_id"]
	db.Where(m).Find(&todos).Count(&count)

	if err != nil {
		log.Println("Error finding record!", err)
		Response.Data["user"] = nil
		Response.Message = "Error finding record"
		Response.Status = "Failed"
		JSON(w, http.StatusNotFound, &Response)
	} else if count < 1 {
		fmt.Println("You do not have any todos")
		Response.Data["user"] = nil
		Response.Message = "You do not have any todos"
		Response.Status = "Failed"
		JSON(w, http.StatusNotFound, &Response)
	} else {
		JSON(w, http.StatusOK, &todos)
	}
	// json.NewEncoder(w).Encode(todos)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	db = db.WithContext(context.TODO())

	m := make(map[string]interface{})
	m["id"] = params["id"]
	m["user_id"] = Response.Data["user_id"]

	todoMap := make(map[string]interface{})
	todoMap["id"] = ""
	todoMap["text"] = ""

	// db.Where(m).Find(&todo)
	// json.NewEncoder(w).Encode(todo)
	var count int64 = 0
	db.Table("todos").Select("id, text").Where(m).Scan(todoMap).Count(&count)
	// db.Raw("SELECT id, text FROM todos WHERE id=? AND user_id=?", params["id"], Response.Data["user_id"]).Scan(todoMap).Count(&count)

	// fmt.Println(todoMap)
	// fmt.Println(count)

	if count < 1 {
		fmt.Println("Todo does not exist")
		// Response.Data["todo"] = nil
		// Response.Message = "Todo does not exist!"
		// Response.Status = "Failed"

		miniResponse := make(map[string]interface{})
		miniResponse["todo"] = nil
		miniResponse["Message"] = "Todo does not exist!"
		miniResponse["Status"] = "Failed"
		JSON(w, http.StatusOK, &miniResponse)
	} else {
		JSON(w, http.StatusOK, &todoMap)
	}
	// json.NewEncoder(w).Encode(todoMap)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// db = dbConn()

	// var todo Todos
	json.NewDecoder(r.Body).Decode(&todo)

	// get signed-in user
	todo.UserID = 1

	// defer db.Close()

	// statement := `INSERT INTO todos(text) VALUES($1) RETURNING id, text`
	// err := db.QueryRow(statement, todo.Text).Scan(&todo.ID, &todo.Text)
	db = db.WithContext(context.TODO())
	err := db.Raw(`INSERT INTO todos(text) VALUES($1) RETURNING id`).Scan(&todo.ID)
	if err != nil {
		fmt.Println("Error when inserting: ", err)
		panic(err)
	}
	fmt.Println(todo)
	// json.NewEncoder(w).Encode(todo)
	JSON(w, http.StatusOK, &todo)
}

func EditTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// var todo Todos
	json.NewDecoder(r.Body).Decode(&todo)

	// db := dbConn()
	// defer db.Close()

	// statement := `UPDATE todos SET text=$2 WHERE id=$1`

	// _, err = db.Exec(statement, params["id"], todo.Text)
	_ = db.Raw(`UPDATE todos SET text=$2 WHERE id=$1`, params["id"], todo.Text)
	if err != nil {
		log.Fatal("Error updating record:", err)
	}

	// statement = `SELECT id, text FROM todos WHERE id=$1`

	// rows := db.QueryRowContext(context.TODO(), statement, params["id"])
	db = db.WithContext(context.TODO())
	db.Select("todos.id, todos.text").Where(params["id"]).Scan(&todo)
	if err != nil {
		log.Fatal("Error getting new record:", err)
	}

	// json.NewEncoder(w).Encode(todo)
	JSON(w, http.StatusOK, &todo)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// db := dbConn()
	// defer db.Close()

	params := mux.Vars(r)

	// statement := `DELETE FROM todos WHERE id=$1`

	db.Raw(`DELETE FROM todos WHERE id=$1`, params["id"])
	if err != nil {
		log.Fatal("Error deleting record:", err)
	}

	// json.NewEncoder(w).Encode("row deleted")
	JSON(w, http.StatusOK, "Deleted Todo")
}

func SearchTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// var tods []Todos
	// var todo Todos

	// db = dbConn()
	fmt.Println("r:", r)
	params := mux.Vars(r)
	fmt.Println(params["key"])
	// statement := `SELECT id, text FROM todos WHERE text LIKE '%' || $1 || '%'`

	// rows, err := db.Query(statement, params["key"])
	rows, err := db.Raw(`SELECT id, text FROM todos WHERE text LIKE '%' || $1 || '%'`, params["key"]).Rows()

	if err != nil {
		fmt.Println("Error when searching: ", err)
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&todo.ID, &todo.Text)
		if err != nil {
			log.Println("Error occured getting todos:", err)
		}
		fmt.Println(todo)
		todos = append(todos, todo)
	}
	// json.NewEncoder(w).Encode(todos)
	JSON(w, http.StatusOK, &todos)
	// defer db.Close()
	defer rows.Close()
}
