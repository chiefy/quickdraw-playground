package quickdraw

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

const (
	qdURL                 = "https://nylottery.ny.gov/quick-draw/past-winning-numbers?page=%d"
	qdRowClass            = "quickdraw-accordion-row"
	qdDateClass           = "result-date"
	qdDrawNumClass        = "draw-number"
	qdWinningNumbersClass = "winning-numbers"
	qdExtraClass          = "col-xs-12 col-sm-1 text-dark-blue"
)

func getDateTime(dn *html.Node) (*time.Time, error) {
	dd, ok := scrape.Find(dn, scrape.ByClass(qdDateClass))
	if !ok {
		return nil, fmt.Errorf("could not decode date in element")
	}
	dt, err := time.Parse("01/02/2006 15:04 PM", scrape.Text(dd))
	if err != nil {
		return nil, err
	}
	return &dt, nil
}

func getDrawNumber(dn *html.Node) (int, error) {
	n, ok := scrape.Find(dn, scrape.ByClass(qdDrawNumClass))
	if !ok {
		return 0, fmt.Errorf("could not decode draw number in element")
	}
	drawNumRaw := strings.Split(scrape.Text(n), "#")
	if len(drawNumRaw) != 2 {
		return 0, fmt.Errorf("could not decode draw number trying to get rid of #")
	}
	return strconv.Atoi(drawNumRaw[1])
}

func getWinningNumbers(dn *html.Node) ([]int64, error) {
	n, ok := scrape.Find(dn, scrape.ByClass(qdWinningNumbersClass))
	if !ok {
		return nil, fmt.Errorf("could not decode winning numbers in element")
	}
	d := scrape.FindAll(n, scrape.ByClass("numbers"))
	if len(d) != 2 {
		return nil, fmt.Errorf("could not decode winning numbers in element")
	}
	wnRaw := append(strings.Split(scrape.Text(d[0]), "-"), strings.Split(scrape.Text(d[1]), "-")...)
	wn := []int64{}
	for _, num := range wnRaw {
		n, err := strconv.Atoi(num)
		if err != nil {
			continue
		}
		wn = append(wn, int64(n))
	}
	return wn, nil
}

func getExtra(dn *html.Node) (int, error) {
	t := strings.Split(scrape.Text(dn), "X")
	r := strings.TrimPrefix(t[2], " ")
	rawNum, err := strconv.Atoi(r)
	if err != nil {
		return 0, err
	}
	return rawNum, nil
}

func scrapeLive() (Draws, error) {
	resp, err := http.Get(fmt.Sprintf(qdURL, 1))
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	draws := []*Draw{}
	drawNodes := scrape.FindAll(root, scrape.ByClass(qdRowClass))
	log.Println(drawNodes)
	if len(drawNodes) == 0 {
		return nil, fmt.Errorf("no html nodes found when scraping")
	}
	for _, d := range drawNodes {
		dt, err := getDateTime(d)
		if err != nil {
			log.Println(err)
			continue
		}
		dn, err := getDrawNumber(d)
		if err != nil {
			log.Println(err)
			continue
		}
		wn, err := getWinningNumbers(d)
		if err != nil {
			log.Println(err)
			continue
		}
		ex, err := getExtra(d.FirstChild.LastChild)
		if err != nil {
			log.Println(err)
			continue
		}
		nd := &Draw{
			DrawNumber:     dn,
			DrawDate:       *dt,
			DrawTime:       fmt.Sprintf("%02d:%02d", dt.Hour(), dt.Minute()),
			WinningNumbers: wn,
			Extra:          ex,
		}
		draws = append(draws, nd)
	}
	return draws, nil
}
