package quickdraw

import (
	"github.com/lib/pq"

	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Draw represents a quickdraw entry
type Draw struct {
	DrawDate       time.Time `json:"draw_date"`
	DrawNumber     int       `json:"draw_number"`
	DrawTime       string    `json:"draw_time"`
	WinningNumbers []int64   `json:"winning_numbers"`
	Extra          int       `json:"extra_multiplier"`
}

// Draws is a collection of Draw
type Draws []*Draw

const (
	fetchURL = "https://data.ny.gov/resource/7sqk-ycpk.json"
	dateForm = "2006-01-02T00:00:00.000"
	// LowPick is the lowest pick number available
	LowPick = 1
	// HighPick is the highest pick number available
	HighPick          = 80
	refreshViewsQuery = `
	REFRESH MATERIALIZED VIEW freq_1day;
	REFRESH MATERIALIZED VIEW freq_7day;
	REFRESH MATERIALIZED VIEW freq_30day;
	REFRESH MATERIALIZED VIEW freq_all_time;
	`
)

// RefreshViews will refresh the materialized views
func RefreshViews(db *sql.DB) error {
	_, err := db.Query(refreshViewsQuery)
	if err != nil {
		return err
	}
	return nil
}

// ScrapeLiveDraws scrapes the live website for new draw data and inserts it into db
func ScrapeLiveDraws(db *sql.DB) error {
	newDraws, err := scrapeLive()
	if err != nil {
		return err
	}
	for _, d := range newDraws {
		err = d.CheckAndInsert(db)
		if err != nil {
			return err
		}
	}
	return RefreshViews(db)
}

// ImportLatest imports the latest entries from the API
func ImportLatest(db *sql.DB) error {
	resp, err := http.Get(fetchURL)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	newDraws := []*Draw{}
	err = json.Unmarshal(body, &newDraws)
	if err != nil {
		return err
	}
	for _, d := range newDraws {
		err = d.CheckAndInsert(db)
		if err != nil {
			return err
		}
	}
	return RefreshViews(db)
}

func getFreqView(viewName string, db *sql.DB) (map[int]int, error) {
	counts := map[int]int{}
	query := fmt.Sprintf("SELECT * FROM %s", viewName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var pick, ct int
		if err := rows.Scan(&pick, &ct); err != nil {
			return nil, err
		}
		counts[pick] = ct
	}
	return counts, nil
}

// GetWinningNumbersFor retrieves a view in the database relating to the past winning numbers based on the viewname
func GetWinningNumbersFor(viewName string, db *sql.DB) (map[int]int, error) {
	switch viewName {
	case "oneday":
		return getFreqView("freq_1day", db)
	case "oneweek":
		return getFreqView("freq_7day", db)
	case "onemonth":
		return getFreqView("freq_30day", db)
	case "alltime":
		return getFreqView("freq_all_time", db)
	}
	return nil, fmt.Errorf("no valid freq time specified")
}

// GetTotalRowsCount gets the total number of records in the database for pagination purposes
func GetTotalRowsCount(db *sql.DB) (int, error) {
	var totalRows int
	totalRowsQuery := "SELECT count(*) FROM draws"
	row := db.QueryRow(totalRowsQuery)
	err := row.Scan(&totalRows)
	if err != nil {
		return 0, err
	}
	return totalRows, nil
}

// GetDraws queries all draws based on the arguments provided
func GetDraws(db *sql.DB, pageNum int, pageSize int, orderBy string, sortDir string) (Draws, error) {
	var offset int = 0
	if pageNum > 1 {
		offset = pageNum * pageSize
	}
	query := fmt.Sprintf("SELECT * FROM draws ORDER BY draws.%s %s OFFSET %d LIMIT %d", orderBy, sortDir, offset, pageSize)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	draws := []*Draw{}
	defer rows.Close()
	for rows.Next() {
		var dn, de int
		var wn pq.Int64Array
		var dd, dt string
		if err := rows.Scan(&dn, &dd, &dt, &wn, &de); err != nil {
			return nil, err
		}
		parsedDate, _ := time.Parse(time.RFC3339, dd)
		d := &Draw{
			DrawNumber:     dn,
			DrawDate:       parsedDate,
			Extra:          de,
			WinningNumbers: []int64(wn),
			DrawTime:       dt,
		}
		draws = append(draws, d)
	}
	return draws, nil
}

// WinningNumbersString converts the winning numbers array to a string value
func (d Draw) WinningNumbersString() string {
	ss := ""
	for i, n := range d.WinningNumbers {
		var comma = ","
		if i == 0 {
			comma = ""
		}
		ss = fmt.Sprintf("%s%s %s", ss, comma, strconv.Itoa(int(n)))
	}
	return ss
}

// Insert inserts a draw into the database
func (d Draw) Insert(db *sql.DB) error {
	a := fmt.Sprintf("ARRAY [%s]", d.WinningNumbersString())
	insert := fmt.Sprintf("INSERT INTO draws (id, drawdate, drawtime, picks, extra) VALUES (%d, '%s', '%s', %s, %d)",
		d.DrawNumber, d.DrawDate.Format("2006-01-02"), d.DrawTime, a, d.Extra)
	_, err := db.Exec(insert)
	if err != nil {
		return err
	}
	log.Println("inserted new record for ", d.DrawNumber)
	return nil
}

// CheckAndInsert checks if a draw exists and if it doesn't, inserts it into the database
func (d Draw) CheckAndInsert(db *sql.DB) error {
	var found int
	row := db.QueryRow(fmt.Sprintf("SELECT id FROM draws WHERE id = %d", d.DrawNumber))
	err := row.Scan(&found)
	switch {
	case err == sql.ErrNoRows:
		return d.Insert(db)
	case err == nil:
		return nil
	default:
		return err
	}
}

// UnmarshalJSON unmarshals custom into Draw struct
func (d *Draw) UnmarshalJSON(b []byte) error {
	var temp map[string]string
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	d.DrawNumber, err = strconv.Atoi(temp["draw_number"])
	if err != nil {
		return err
	}
	d.Extra, err = strconv.Atoi(temp["extra_multiplier"])
	if err != nil {
		return err
	}
	d.DrawDate, err = time.Parse(dateForm, temp["draw_date"])
	if err != nil {
		return err
	}
	d.DrawTime = temp["draw_time"]
	if err != nil {
		return err
	}
	winNums := strings.Split(temp["winning_numbers"], " ")
	d.WinningNumbers = make([]int64, 0)
	for _, n := range winNums {
		w, err := strconv.Atoi(n)
		if err != nil {
			return err
		}
		d.WinningNumbers = append(d.WinningNumbers, int64(w))
	}
	return nil
}
