package main

import (
	"database/sql"
	"github.com/chiefy/quick-draw-explorer/pkg"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"fmt"
	"log"
	"os"
)

var (
	doImport = kingpin.Flag("import", "Import and parse CSV from the internet").Bool()
	doPoll   = kingpin.Flag("poll", "Poll the API for latest data").Bool()
)

const (
	csvFile  = "quick-draw.csv"
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "supersecret"
	dbname   = "quickdraw"
)

func connectToDb() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return quickdraw.GetDb(connString)
}

func main() {
	kingpin.Parse()

	if *doImport {
		db := connectToDb()
		defer db.Close()
		err := quickdraw.ImportPicks(db, csvFile)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	if *doPoll {
		db := connectToDb()
		defer db.Close()
		err := quickdraw.ImportLatest(db)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
	//log.Printf("%#v\n", *doImport)

}
