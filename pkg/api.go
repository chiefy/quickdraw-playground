package quickdraw

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	pageSize = 50
)

// API represents the API singleton
type API struct {
	Db          *sql.DB
	Host        string
	Port        int
	AllowedURLs []string
}

type TableResponse struct {
	Draws    []*Draw `json:"draws"`
	Total    int     `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func (a API) listDraws(w http.ResponseWriter, r *http.Request) {
	pageNum, err := strconv.Atoi(chi.URLParam(r, "pageNum"))
	if err != nil {
		log.Println("could not convert p url param to int, using 1", err)
		pageNum = 1
	}
	pageSize, err := strconv.Atoi(chi.URLParam(r, "pageSize"))
	if err != nil {
		log.Println("could not convert rpp url param to int, using 25")
		pageSize = 25
	}
	sortBy := chi.URLParam(r, "sortBy")
	sortDir := chi.URLParam(r, "sortDir")

	draws, err := GetDraws(a.Db, pageNum, pageSize, sortBy, sortDir)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("GetDraws error %s", err)))
		return
	}

	totalCount, err := GetTotalRowsCount(a.Db)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("GetTotalRowsCount error %s", err)))
		return
	}
	res := &TableResponse{
		Draws:    draws,
		Total:    totalCount,
		Page:     pageNum,
		PageSize: pageSize,
	}

	drawsJSON, err := json.Marshal(&res)
	if err != nil {
		w.Write([]byte("JSON marshaling error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(drawsJSON)
}

func (a API) freqView(w http.ResponseWriter, r *http.Request) {
	viewName := chi.URLParam(r, "viewname")
	counts, err := GetWinningNumbersFor(viewName, a.Db)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	d, err := json.Marshal(counts)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

// Serve starts the API server
func (a *API) Serve() {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   a.AllowedURLs,
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler)

	r.Route("/api", func(r chi.Router) {
		// r.With(paginate).Get("/draws", q.listDraws)
		r.Get("/draws/p{pageNum}/s{pageSize}/b{sortBy}/{sortDir}", a.listDraws)
		r.Get("/freq/{viewname}", a.freqView)
	})

	s := fmt.Sprintf("%s:%d", a.Host, a.Port)
	log.Println("Quickdraw Explorer Serving on " + s)
	log.Fatal(http.ListenAndServe(s, r))
}
