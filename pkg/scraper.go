package quickdraw

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

const (
	qdURL      = "https://nylottery.ny.gov/quick-draw/past-winning-numbers?page=%d"
	qdRowClass = "quickdraw-accordion-row"
)

type DrawResults struct {
	rawDraws []*html.Node
}

func Scrape() (Draws, error) {
	resp, err := http.Get(fmt.Sprintf(qdURL, 1))
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	draws := scrape.FindAll(root, scrape.ByClass(qdRowClass))
	for _, d := range draws {
		dd, ok := scrape.Find(d, scrape.ByClass("result-date"))
		if !ok {
			continue
		}
		dt, err := time.Parse("01/02/2006 15:04 PM", scrape.Text(dd))
		if err != nil {
			log.Println("Could not parse elements date/time", err)
			continue
		}

		log.Println(dt)
	}

	return []*Draw{}, nil
}
