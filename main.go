package main

import (
	"log"
	"time"
)

func main() {
	//err := importPicks()
	end := time.Now()
	diff, _ := time.ParseDuration("120h")
	start := end.Add(-diff)
	dq := NewDrawQuery(
		&start,
		&end,
	)
	db, err := getDb()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = dq.Exec(db)
	if err != nil {
		log.Fatal(err)
	}

}
