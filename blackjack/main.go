package main

import (
	"flag"
	"log"
)

var numPlayers = 2

func init() {
	flag.IntVar(&numPlayers, "num-players", numPlayers, "number of players (min 2 - max 7)")
	flag.Parse()

	if numPlayers > 7 || numPlayers < 2 {
		log.Fatal("number of players must be within the range of 2 - 7")
	}
}

func main() {
	newGame(numPlayers).Play()
}
