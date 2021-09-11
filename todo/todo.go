package todo

import "database/sql"

type Todos struct {
	ID   string `json:"id"`
	Item string `json:"item"`
}

var (
	db  *sql.DB
	err error
)
