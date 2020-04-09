package main

import (
	"fmt"
	"log"

	"github.com/fuskovic/gophercises/normalizer/store"
)

func main() {
	s, err := store.New()
	if err != nil {
		log.Fatalf("failed to create new store instance : %s\n", err)
	}

	if err := s.Normalize(); err != nil {
		log.Fatalf("failed to normalize database : %s\n", err)
	}

	normalizedEntries, err := s.List()
	if err != nil {
		log.Fatalf("failed to list numbers : %s\n", err)
	}

	for _, e := range normalizedEntries {
		fmt.Println(e.PhoneNumber)
	}
}
