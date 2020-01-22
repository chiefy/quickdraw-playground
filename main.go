package main

import (
	"github.com/chiefy/quick-draw-explorer/pkg"
	"github.com/jasonlvhit/gocron"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"strconv"

	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	doImport    = kingpin.Flag("import", "Import and parse CSV from the internet").Bool()
	doPoll      = kingpin.Flag("poll", "Poll the API for latest data").Bool()
	doServe     = kingpin.Flag("serve", "Start the API server").Bool()
	allowedURLs string

	apiPort int
	apiHost string
)

const (
	csvFile    = "quick-draw.csv"
	dbHost     = "db"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "supersecret"
	dbName     = "quickdraw"

	defaultAPIPort    = 9090
	defaultAPIHost    = "localhost"
	defaultAllowedURL = "http://localhost:8080"
)

func init() {
	allowedURLs = os.Getenv("ALLOWED_CORS_URLS")
	if allowedURLs == "" {
		allowedURLs = defaultAllowedURL
	}
	var err error
	apiPortStr := os.Getenv("API_PORT")
	if apiPort, err = strconv.Atoi(apiPortStr); err != nil {
		apiPort = defaultAPIPort
	}

	apiHost = os.Getenv("API_HOST")
	if apiHost == "" {
		apiHost = defaultAPIHost
	}
}

func doFetchLatest() {
	db := connectToDb()
	defer db.Close()
	err := quickdraw.ImportLatest(db)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchLatestDrawsCron() {
	gocron.Every(15).Minutes().Do(doFetchLatest)
	<-gocron.Start()
}

func connectToDb() *sql.DB {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	return quickdraw.GetDb(connString)
}

func main() {
	kingpin.Parse()

	if *doImport {
		db := connectToDb()
		defer db.Close()
		for {
			log.Println("Checking for connection to DB...")
			err := db.Ping()
			if err != nil {
				log.Println("Could not connect to DB, retrying in 5s...")
			} else {
				break
			}
			time.Sleep(5 * time.Second)
		}
		err := quickdraw.ImportPicks(db, csvFile)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	if *doPoll {
		doFetchLatest()
		os.Exit(0)
	}
	if *doServe {
		db := connectToDb()
		defer db.Close()

		go fetchLatestDrawsCron()
		a := quickdraw.API{
			Db:          db,
			Host:        apiHost,
			Port:        apiPort,
			AllowedURLs: []string{allowedURLs},
		}
		a.Serve()
	}
}
