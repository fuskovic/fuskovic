package service

import (
	"context"
	"net/http"

	"github.com/fuskovic/gophercises/urlshortener/handler"
	"github.com/fuskovic/gophercises/urlshortener/store"
	"github.com/gorilla/mux"
)

const notFound = "https://9gag.com/404"

var fb fallback

type (
	// Service represents the url shortening service
	Service struct {
		Router *mux.Router
		Store  *store.Store
	}

	fallback http.HandlerFunc
)

// ServeHTTP is used by the fallback handlerFunc to handle 404s
func (fb fallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, notFound, http.StatusPermanentRedirect)
}

// New creates an instance of this service
func New(ctx context.Context, s *store.Store, m map[string]string) Service {
	mapHandler := handler.ByteDataHandler(m, fb)
	dbHandler := handler.DbHandler(ctx, *s, fb)
	r := mux.NewRouter()
	r.HandleFunc("/linkedin", mapHandler)
	r.HandleFunc("/github", mapHandler)
	r.HandleFunc("/db-linkedin", dbHandler)
	r.HandleFunc("/db-github", dbHandler)
	return Service{
		Router: r,
		Store:  s,
	}
}
