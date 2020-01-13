package main

import (
	"github.com/chiefy/quick-draw-explorer/pkg"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	doImport = kingpin.Flag("import", "Import and parse CSV from the internet").Bool()
	doPoll   = kingpin.Flag("poll", "Poll the API for latest data").Bool()
	doServe  = kingpin.Flag("serve", "Start the API server").Bool()
)

const (
	csvFile  = "quick-draw.csv"
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "supersecret"
	dbname   = "quickdraw"

	apiPort = 9090
	apiHost = "localhost"
)

func connectToDb() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return quickdraw.GetDb(connString)
}

func doServeAPI() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		db := connectToDb()
		defer db.Close()
		counts, err := quickdraw.GetWinningNumbersCount(db)
		if err != nil {
			w.Write([]byte("error"))
			return
		}
		log.Println(counts)
		d, err := json.Marshal(counts)
		if err != nil {
			w.Write([]byte("error"))
			return
		}
		w.Write(d)
	})
	s := fmt.Sprintf("%s:%d", apiHost, apiPort)
	log.Println("Quickdraw Explorer Serving on " + s)
	log.Fatal(http.ListenAndServe(s, r))
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
	if *doServe {
		doServeAPI()
	}
}
