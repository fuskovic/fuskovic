package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fuskovic/gophercises/urlshortener/handler"
	"github.com/fuskovic/gophercises/urlshortener/service"
	"github.com/fuskovic/gophercises/urlshortener/store"

	"gopkg.in/yaml.v2"
)

const (
	addr   = ":8080"
	sqlite = "sqlite3"
)

var (
	startMsg string
	config   = "mappings.yaml"
)

func init() {
	flag.StringVar(&config, "config", config, "path to yaml or json config file")
	flag.Parse()
	startMsg = fmt.Sprintf("url shortener started on %s\nconfig used : %s\n", addr, config)
}

func isYaml() bool {
	ext := filepath.Ext(config)
	return ext == ".yml" || ext == ".yaml"
}

func getBinPath() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(wd, "store/db.sqlite3")
}

func parse(data []byte, isYaml bool) (map[string]string, error) {
	var o handler.Object
	var err error

	if isYaml {
		err = yaml.Unmarshal(data, &o.Routes)
	} else {
		err = json.Unmarshal(data, &o)
	}

	if err != nil {
		return o.AsMap, err
	}
	return o.Map(), nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	data, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("failed to read %s : %s\n", config, err.Error())
	}

	m, err := parse(data, isYaml())
	if err != nil {
		log.Fatalf("failed to parse data : %s\n", err.Error())
	}

	db, err := sql.Open(sqlite, getBinPath())
	if err != nil {
		log.Fatalf("failed to connect to sqlite : %s\n", err.Error())
	}
	defer db.Close()

	s, err := store.New(ctx, db, m)
	if err != nil {
		log.Fatalf("failed to initialize store : %s\n", err.Error())
	}

	svc := service.New(ctx, s, m)
	println(startMsg)
	http.ListenAndServe(addr, svc.Router)
}
