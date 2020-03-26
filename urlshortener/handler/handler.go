package handler

import (
	"context"
	"net/http"

	"github.com/fuskovic/gophercises/urlshortener/store"
)

type (
	// Route stores a path and it's respective URL
	Route struct {
		Path string
		URL  string
	}
	// Object is a wrapper for yaml and json data
	Object struct {
		Routes []Route `yaml:"paths" json:"paths"`
		AsMap  map[string]string
	}
)

// Map converts routes to a map
func (o *Object) Map() map[string]string {
	o.AsMap = make(map[string]string)
	for _, route := range o.Routes {
		o.AsMap[route.Path] = route.URL
	}
	return o.AsMap
}

// MapHandler handles requests requests by finding urls for a given path
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.String()]
		if !ok {
			fallback.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
		}
	}
}

// ByteDataHandler is a wrapper for the map handler
func ByteDataHandler(m map[string]string, fallback http.Handler) http.HandlerFunc {
	return MapHandler(m, fallback)
}

// DbHandler is a wrapper for the stores handler method
func DbHandler(ctx context.Context, s store.Store, fallback http.Handler) http.HandlerFunc {
	return s.Handler(ctx, fallback)
}
