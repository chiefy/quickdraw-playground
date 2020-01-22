package quickdraw

import (
	// use postgres driver
	_ "github.com/lib/pq"

	"database/sql"
)

const (
	schemaName = "quickdraw"
)

// GetDb opens the database w/ postgres driver
func GetDb(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}
