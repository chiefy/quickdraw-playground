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

type Draws []*Draw

const (
	fetchURL = "https://data.ny.gov/resource/7sqk-ycpk.json"
	dateForm = "2006-01-02T00:00:00.000"

	LowPick  = 1
	HighPick = 80
)

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
	return nil
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

func GetDraws(db *sql.DB, pageSize int) (Draws, error) {
	query := fmt.Sprintf("SELECT * FROM draws ORDER BY draws.id DESC LIMIT %d", pageSize)
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
