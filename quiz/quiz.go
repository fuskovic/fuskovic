package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	path    = "./problems.csv"
	seconds = 30
)

func init() {
	flag.StringVar(&path, "path", path, "path to csv file with questions")
	flag.IntVar(&seconds, "seconds", seconds, "time limit for quiz in seconds")
	flag.Parse()
}

func buildQuiz(path string) (map[string]int, error) {
	quiz := make(map[string]int)

	file, err := os.Open(path)
	if err != nil {
		return quiz, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return quiz, err
	}

	for _, record := range records {
		answer, err := strconv.Atoi(record[1])
		if err != nil {
			return quiz, err
		}
		quiz[record[0]] = answer
	}
	return quiz, nil
}

func clean(s string) string {
	return strings.Trim(s, " ")
}

func printResult(correct, numQuestions int) {
	fmt.Printf("times up!\n%d/%d correct\n", correct, numQuestions)
}

func main() {
	quiz, err := buildQuiz(path)
	if err != nil {
		log.Fatalf("failed to build quiz : %s", err.Error())
	}

	var correct int
	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	scanner := bufio.NewScanner(os.Stdin)

	go func(c *int) {
		<-timer.C
		printResult(*c, len(quiz))
		os.Exit(1)
	}(&correct)

	for question, answer := range quiz {
		println(question)
		for scanner.Scan() {
			input, err := strconv.Atoi(clean(scanner.Text()))
			if err != nil {
				log.Fatalf("error converting input to int - %s", err.Error())
			}
			if input == answer {
				correct++
			}
			break
		}
	}
	printResult(correct, len(quiz))
}
