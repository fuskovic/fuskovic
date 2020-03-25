package handler

import (
	"net/http"

	"gopkg.in/yaml.v2"
)

// Route stores a path and it's respective URL
type Route struct {
	Path string
	URL  string
}

// YamlObject is a wrapper for the yaml data
type YamlObject struct {
	Routes []Route `yaml:"paths"`
	AsMap  map[string]string
}

// Map converts routes to a map
func (y *YamlObject) Map() map[string]string {
	y.AsMap = make(map[string]string)
	for _, route := range y.Routes {
		y.AsMap[route.Path] = route.URL
	}
	return y.AsMap
}

func parseYAML(yml []byte) (map[string]string, error) {
	var y YamlObject
	if err := yaml.Unmarshal(yml, &y.Routes); err != nil {
		return y.AsMap, err
	}
	return y.Map(), nil
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

// YAMLHandler handles requests by checking the yaml config for a url belonging to a path.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	fb := func(w http.ResponseWriter, r *http.Request) { fallback.ServeHTTP(w, r) }

	m, err := parseYAML(yml)
	if err != nil {
		return fb, err
	}
	return MapHandler(m, fallback), nil
}
