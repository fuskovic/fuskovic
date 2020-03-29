package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	arcToStartOn = "intro"
	adventureCfg = "gopher.json"
	toStdOut     bool
	chapterFiles []string
)

func init() {
	flag.StringVar(&adventureCfg, "story-config", adventureCfg, "path to json config for the adventure")
	flag.BoolVar(&toStdOut, "to-std-out", false, "print adventure to stdout vs the browser")
	flag.StringVar(&arcToStartOn, "start-on", arcToStartOn, "choose which part of the story you want to start on")
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

func (a Adventure) run() {
	currentArc, ok := a[arcToStartOn]
	if !ok {
		log.Fatalf("arc to start on : %s does not exist", arcToStartOn)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for currentArc.Title != a["home"].Title {
		fmt.Println()
		for _, line := range currentArc.Stories {
			fmt.Println(line)
		}

		fmt.Println(`
		select an option using the option number
		`)

		for i, option := range currentArc.Options {
			opt := fmt.Sprintf("option number %d : %s\n", i, option.Text)
			fmt.Println(opt)
		}

		for scanner.Scan() {
			input, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Printf("failed to evaluate input : %s\n", err.Error())
			}
			currentArc = a[currentArc.Options[input].Arc]
			break
		}
	}
	fmt.Println("adventure complete!")
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

	if !toStdOut {
		fmt.Printf("starting adventure on %s\n", addr)
		http.ListenAndServe(addr, buildAdventureRouter(adventure))
	} else {
		adventure.run()
	}
}
