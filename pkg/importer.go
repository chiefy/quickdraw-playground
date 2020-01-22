package quickdraw

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	importRoot      = "INSERT INTO draws (id, drawdate, drawtime, picks, extra)\nVALUES\n"
	bulkInsertCount = 1000
)

// ImportPicks imports picks from a CSV
func ImportPicks(db *sql.DB, csvFile string) error {
	csvFileHandle, err := os.Open(fmt.Sprintf("./%s", csvFile))
	if err != nil {
		return err
	}
	r := csv.NewReader(csvFileHandle)

	// disregard first line of field names
	_, err = r.Read()
	if err != nil {
		return err
	}

	sqlStatement := importRoot
	var count = 0
	var totalInserts = 0

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if record[0] == "" {
			log.Println("no date found for draw " + record[1])
			continue
		}
		if count == 0 {
			sqlStatement = importRoot
		} else {
			sqlStatement = fmt.Sprintln(sqlStatement + ",")
		}
		_, err = strconv.Atoi(record[1])
		if err != nil {
			return err
		}
		_, err = strconv.Atoi(record[4])
		if err != nil {
			return err
		}
		splitPicks := strings.Split(record[3], " ")
		joinedPicks := fmt.Sprintf("ARRAY [%s]", strings.Join(splitPicks, ","))
		sqlStatement = fmt.Sprint(sqlStatement + "(" + record[1] + "," + "'" + record[0] +
			"','" + record[2] + "'," + joinedPicks + "," + record[4] + ")")

		if count == bulkInsertCount {
			_, err = db.Exec(sqlStatement)
			if err != nil {
				log.Println(sqlStatement)
				log.Println(err)
				os.Exit(1)
			}
			totalInserts += count
			log.Println("Bulk inserting", count, totalInserts)
			count = 0
		} else {
			count++
		}
	}

	_, err = db.Exec(sqlStatement)
	if err != nil {
		return err
	}

	return RefreshViews(db)
}
