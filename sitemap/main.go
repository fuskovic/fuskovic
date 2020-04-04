package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/fuskovic/gophercises/sitemap/crawler"
)

var (
	baseURL = "https://gophercises.com"
	depth   = 0
)

func init() {
	flag.StringVar(&baseURL, "baseURL", baseURL, "base url for which to build sitemap for")
	flag.IntVar(&depth, "depth", depth, "number of levels to check for links")
	flag.Parse()
}

func main() {
	URL, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("%s is an invalid url : %s\n", baseURL, err)
	}

	fmt.Printf("building sitemap for %v\n", URL)

	c, err := crawler.Init(baseURL, depth)
	if err != nil {
		log.Fatalf("failed to initialize crawler : %s\n", err)
	}

	var links []string

	if depth > 0 {
		links = c.CrawlByDepth()
	} else {
		links = c.CrawlAll()
	}

	xmlData, err := crawler.ToXML(links)
	if err != nil {
		log.Fatalf("failed to convert links to XML : %s\n", err)
	}

	if err := ioutil.WriteFile("sitemap.xml", xmlData, 0666); err != nil {
		log.Fatalf("failed to write file : %s\n", err)
	}

	fmt.Println("sitemap.xml successfully created")
}
