package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fuskovic/gophercises/deck"
)

const (
	none = iota
	stand
	hit
)

type Game struct {
	Deck    *deck.Deck
	Players *Players
}

func newGame(playersToAdd int) *Game {
	var players Players
	dealer := Player{Name: "dealer"}

	for i := 1; i < playersToAdd; i++ {
		newPlayer := Player{Name: fmt.Sprintf("player%d", i)}
		players = append(players, newPlayer)
	}

	players = append(players, dealer)

	d := deck.New()
	var nextCard deck.Card

	fmt.Printf("***************STARTING NEW GAME WITH %d PLAYERS***************\n", len(players))

	// give each player two cards
	for i := 0; i < 2; i++ {
		for i, p := range players {
			nextCard, d = d.Draw()
			players[i] = p.Hit(nextCard)
		}
	}

	return &Game{
		Deck:    d,
		Players: &players,
	}
}

func (g *Game) GetUserInput() (selection int) {
	scanner := bufio.NewScanner(os.Stdin)

	isInvalidInput := func(input int, err error) bool {
		return input > hit || input < stand || err != nil
	}

	for scanner.Scan() {
		input, err := strconv.Atoi(scanner.Text())
		if isInvalidInput(input, err) {
			println("invalid selection - enter 1 for stand or 2 for hit")
			continue
		}
		selection = input
		break
	}
	return
}

func (g *Game) GetWinner() (winner Player) {
	hasBestScore := func(p Player, scoreToBeat uint) bool {
		return p.Score < 22 && p.Score > scoreToBeat
	}

	for _, p := range *g.Players {
		if hasBestScore(p, winner.Score) {
			winner = p
		}
	}
	return
}

func (g *Game) Status() {
	for _, player := range *g.Players {
		if player.Name == "player1" || player.Name == "dealer" {
			fmt.Printf("%s has...\n", player.Name)
			for _, card := range player.Hand {
				fmt.Printf("%s of %s\n", card.Value, card.Suit)
			}
			fmt.Printf("%s score : %d\n", player.Name, player.Score)
		}
	}
}

func handleBustedPlayer(p Player) {
	if p.Name == "dealer" {
		log.Fatal("dealer has busted - YOU WIN!")
	}

	if p.Name == "player1" {
		log.Fatal("dealer wins")
	}

	fmt.Printf("%s has busted\n", p.Name)
}

func (g *Game) Play() {
	selection := none
	var nextCard deck.Card

	players := *g.Players
	g.Status()

	for i, p := range players {
		for selection == none || selection == hit {
			if p.Name == "player1" {
				fmt.Println("it's your turn\n1 - stay\n2 - hit")
				selection = g.GetUserInput()
			} else {
				selection = p.getAutomatedInput()
			}

			if selection == hit {
				fmt.Printf("%s has chose to hit\n", p.Name)
				nextCard, g.Deck = g.Deck.Draw()
				players[i] = p.Hit(nextCard)
				if p.Busted() {
					handleBustedPlayer(players[i])
				}
				g.Status()
			} else {
				fmt.Printf("%s has chose to stay\n", p.Name)
			}
		}
		selection = none
	}
	g.Players = &players
	winner := g.GetWinner()
	fmt.Printf("the best hand belongs to %s with a high score of %d\n", winner.Name, winner.Score)
}
