package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fuskovic/gophercises/linkparser/linx"
)

var htmlFiles = []string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}

func main() {
	for _, f := range htmlFiles {
		data, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("failed to read %s : %s\n", f, err)
		}
		links, err := linx.GetLinks(data)
		if err != nil {
			log.Fatalf("failed to get links : %s\n", err.Error())
		}
		fmt.Printf("found links in %s\n", f)
		for _, link := range links {
			fmt.Printf("\tthe text '%s' goes to the link '%s'\n", link.Text, link.Href)
		}
	}
}
