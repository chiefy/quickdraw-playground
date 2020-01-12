package quickdraw

import (
	_ "github.com/lib/pq"

	"database/sql"
)

const (
	schemaName = "quickdraw"
)

func GetDb(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
