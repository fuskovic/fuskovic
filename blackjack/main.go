package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	numPlayers = 2
	rounds     = 1
)

func init() {
	flag.IntVar(&numPlayers, "num-players", numPlayers, "number of players (default : 2, min : 2,  max : 7)")
	flag.IntVar(&rounds, "rounds", rounds, "rounds of blackjack to play ( default : 1, min : 1, max : 10 )")
	flag.Parse()

	if numPlayers > 7 || numPlayers < 2 {
		log.Fatal("number of players must be within the range of 2 - 7")
	}

	if rounds > 10 || rounds < 1 {
		log.Fatal("number of rounds must be within the range 1-10")
	}
}

func main() {
	fmt.Printf("***************STARTING NEW GAME WITH %d PLAYERS***************\n", numPlayers)
	game := NewGame(numPlayers)

	for i := 1; i <= rounds; i++ {
		fmt.Printf("********** STARTING ROUND %d **********\n", i)
		winners := game.Play()
		if len(winners) == 1 {
			fmt.Printf("the winner is %s\n", winners[0].name)
		} else {
			fmt.Println("the winners are...")
			for _, winner := range winners {
				fmt.Println(winner.name)
			}
		}
		fmt.Printf("********** ROUND %d COMPLETE **********\n", i)
		game.Reset()
	}
}
