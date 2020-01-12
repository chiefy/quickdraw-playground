package quickdraw

import (
	_ "github.com/lib/pq"

	"bytes"
	"database/sql"
	"fmt"
	_ "log"
	"sort"
	"text/template"
	"time"
)

const (
	baseQueryTmpl = `
SELECT COUNT(*) AS r FROM (
SELECT Draws.ID, DATE(Draws.drawDate) AS dd 
FROM Draws 
INNER JOIN Picks ON Draws.ID = Picks.drawID 
WHERE Picks.pickNum = {{.PickNum}} AND dd BETWEEN '{{.StartDatePretty}}' AND '{{.EndDatePretty}}'
GROUP BY Draws.ID
) AS f;`
	queryTmplName = "queryTemplate"
	lowPick       = 1
	highPick      = 80
)

type PickFreq struct {
	Num       int
	Frequency int
}

type DrawQueryResult map[int]int

type DrawQuery struct {
	PickNum         int
	StartDateTime   *time.Time
	EndDateTime     *time.Time
	StartDatePretty string
	EndDatePretty   string
	Results         DrawQueryResult
	PickFreq        []*PickFreq
	queryTmpl       *template.Template
}

func NewDrawQuery(start, end *time.Time) *DrawQuery {
	nq := &DrawQuery{
		PickNum:         1,
		StartDateTime:   start,
		EndDateTime:     end,
		StartDatePretty: formattedDate(start),
		EndDatePretty:   formattedDate(end),
		queryTmpl:       template.Must(template.New(queryTmplName).Parse(baseQueryTmpl)),
	}
	return nq
}

func formattedDate(d *time.Time) string {
	return d.Format("2006-01-02")
}

func (dq *DrawQuery) Exec(db *sql.DB) error {
	dq.Results = map[int]int{}

	for pn := lowPick; pn <= highPick; pn++ {
		dq.PickNum = pn
		var r *int
		queryBuf := &bytes.Buffer{}
		err := dq.queryTmpl.Execute(queryBuf, dq)
		if err != nil {
			return err
		}
		rows, err := db.Query(queryBuf.String())
		if err != nil {
			return err
		}
		defer rows.Close()
		if !rows.Next() {
			return fmt.Errorf("No rows for %d", dq.PickNum)
		}
		err = rows.Scan(&r)
		dq.Results[dq.PickNum] = *r
		if err != nil {
			return err
		}

		//log.Printf("results for %d = %d", dq.PickNum, dq.Results[dq.PickNum])
	}

	var pf []*PickFreq
	for k, v := range dq.Results {
		pf = append(pf, &PickFreq{k, v})
	}
	sort.Slice(pf, func(i, j int) bool {
		return pf[i].Frequency > pf[j].Frequency
	})
	for _, kv := range pf {
		fmt.Printf("%d, %d\n", kv.Num, kv.Frequency)
	}
	return nil
}
