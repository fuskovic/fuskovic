package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

const addr = ":8080"

type (
	// Option contains the text and the result of selecting the option.
	Option struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	}

	// Options contains multiple options.
	Options []Option

	// Story represents the sentences in a story
	Story []string

	// Arc represents a single part of the adventure
	Arc struct {
		Title   string  `json:"title"`
		Stories Story   `json:"story"`
		Options Options `json:"options"`
	}

	// Adventure represents all of the sections of the story
	Adventure map[string]Arc
)

var (
	adventureCfg = "gopher.json"
	chapterFiles []string
)

func init() {
	flag.StringVar(&adventureCfg, "story-config", adventureCfg, "path to json config for the adventure")
	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working director : %s\n", err.Error())
	}

	pagesDir := filepath.Join(wd, "pages")
	chapterFiles = append(chapterFiles,
		filepath.Join(pagesDir, "base.html"),
		filepath.Join(pagesDir, "index.html"),
		filepath.Join(pagesDir, "chapter.html"),
	)
}

func chapterHandler(a Arc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("base").ParseFiles(chapterFiles...)
		if err != nil {
			fmt.Printf("failed to build template : %s\n", err.Error())
		}
		if err := tmpl.ExecuteTemplate(w, "base", a); err != nil {
			fmt.Printf("failed to execute template : %s\n", err.Error())
		}
	}
}

func buildAdventureRouter(a Adventure) *mux.Router {
	r := mux.NewRouter()

	for name, chapter := range a {
		var path string
		if name == "intro" {
			path = "/"
		} else {
			path = fmt.Sprintf("/%s", name)
		}
		r.Handle(path, chapterHandler(chapter))
	}

	return r
}

func main() {
	var adventure Adventure

	data, err := ioutil.ReadFile(adventureCfg)
	if err != nil {
		log.Fatalf("failed to read file\nfile : %s\nerror : %s\n", adventureCfg, err.Error())
	}

	if err := json.Unmarshal(data, &adventure); err != nil {
		log.Fatalf("failed to unmarshal json : %s\n", err.Error())
	}

	fmt.Printf("starting adventure on %s\n", addr)
	http.ListenAndServe(addr, buildAdventureRouter(adventure))
}
