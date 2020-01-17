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
)

const (
	pageSize = 50
)

type QuickdrawAPI struct {
	Db          *sql.DB
	Host        string
	Port        int
	AllowedURLs []string
}

func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

func (q QuickdrawAPI) listDraws(w http.ResponseWriter, r *http.Request) {
	draws, err := GetDraws(q.Db, pageSize)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("GetDraws error %s", err)))
		return
	}
	drawsJSON, err := json.Marshal(&draws)
	if err != nil {
		w.Write([]byte("JSON marshaling error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(drawsJSON)
}

func (q QuickdrawAPI) freqView(w http.ResponseWriter, r *http.Request) {
	viewName := chi.URLParam(r, "viewname")
	counts, err := GetWinningNumbersFor(viewName, q.Db)
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

func (q *QuickdrawAPI) Serve() {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   q.AllowedURLs,
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

	r.Route("/", func(r chi.Router) {
		// r.With(paginate).Get("/draws", q.listDraws)
		r.Get("/draws", q.listDraws)
		r.Get("/freq/{viewname}", q.freqView)
	})

	s := fmt.Sprintf("%s:%d", q.Host, q.Port)
	log.Println("Quickdraw Explorer Serving on " + s)
	log.Fatal(http.ListenAndServe(s, r))
}
