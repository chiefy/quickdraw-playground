package main

import (
	_ "github.com/mattn/go-sqlite3"

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
	dbFile  = "quick-draw.db"
	csvFile = "quick-draw.csv"
)

var (
	pickCols = []int{4, 24}
)

func padStr(dateStr string) string {
	num, err := strconv.Atoi(dateStr)
	if err != nil {
		return ""
	}
	if num < 10 {
		return fmt.Sprintf("%d%d", 0, num)
	}
	return fmt.Sprintf("%d", num)
}

func fixDateStr(dateStr string) string {
	splitStr := strings.Split(dateStr, "/")
	if len(splitStr) < 3 {
		log.Printf("Bad date string %s\n", dateStr)
		return "2000-01-01"
	}
	return fmt.Sprintf("20%s-%s-%s", padStr(splitStr[2]), padStr(splitStr[0]), padStr(splitStr[1]))
}

func getDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s", dbFile))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func importPicks() error {
	db, err := getDb()
	defer db.Close()

	if err != nil {
		return err
	}

	csvFile, err := os.Open(fmt.Sprintf("./%s", csvFile))
	if err != nil {
		return err
	}

	r := csv.NewReader(csvFile)
	firstLine := true
	for {
		record, err := r.Read()
		if firstLine {
			firstLine = false
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		drawID, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println(record)
			return err
		}

		if len(record) < 24 {
			fmt.Printf("ERROR! Record has insufficient columns - %#v", record)
			continue
		}

		insertDrawStatement := "insert into Draws(drawDate,ID,drawTime,extraPick) values(?,?,?,?)"
		drawDate := fixDateStr(record[0])
		log.Printf("Inserting draw %d \n", drawID)
		_, err = db.Exec(insertDrawStatement, drawDate, drawID, record[2], record[3])
		if err != nil {
			return err
		}
		insertPickStatement := "insert into Picks(drawID,pickNum) values"
		for colNum := pickCols[0]; colNum < pickCols[1]; colNum++ {
			comma := ""
			newPick, err := strconv.Atoi(record[colNum])
			if err != nil {
				fmt.Println(record)
				return err
			}
			if colNum != pickCols[1]-1 {
				comma = ","
			}
			insertPickStatement = fmt.Sprintf("%s (%d,%d)%s ", insertPickStatement, drawID, newPick, comma)
		}

		_, err = db.Exec(insertPickStatement)
		if err != nil {
			return err
		}

	}
	return nil
}
