package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fuskovic/gophercises/urlshortener/handler"
)

const (
	addr     = ":8080"
	notFound = "https://9gag.com/404"
	config   = "mappings.yaml"
)

var fb fallback

type fallback http.HandlerFunc

func (fb fallback) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, notFound, http.StatusPermanentRedirect)
}

func main() {
	yaml, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("failed to read %s : %s", config, err.Error())
	}

	yh, err := handler.YAMLHandler([]byte(yaml), fb)
	if err != nil {
		fmt.Printf("failed to init yaml handler : %s\n", err.Error())
		os.Exit(1)
	}

	http.HandleFunc("/linkedin", yh)
	http.HandleFunc("/github", yh)
	http.ListenAndServe(addr, nil)
}
