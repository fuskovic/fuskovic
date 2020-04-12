package main

import (
	"fmt"
	"log"

	"github.com/fuskovic/gophercises/deck"
)

type (
	Player struct {
		Name  string
		Score uint
		Hand  []deck.Card
	}
	Players []Player
)

func (p *Player) getAutomatedInput() int {
	if p.Score <= 16 || p.Score == 17 && p.HasAce() {
		return hit
	}
	return stand
}

func (p *Player) Hit(newCard deck.Card) Player {
	p.Hand = append(p.Hand, newCard)

	handleAce := func() {
		if p.Score < 22 && (p.Score+11) < 22 {
			p.Score += 11
		} else if p.Score < 22 && (p.Score+11) > 21 {
			p.Score += 1
		}
	}

	fmt.Printf("%s hit a %s of %s\n", p.Name, newCard.Value, newCard.Suit)

	if newCard.Value == "Ace" {
		handleAce()
	} else {
		p.Score += uint(deck.CardMap[newCard.Value])
	}

	if p.Score == 21 {
		log.Fatalf("%s got 21, %s WINS!", p.Name, p.Name)
	}

	return *p
}

// HasAce evaluates whether or not the players hand contains an Ace.
func (p *Player) HasAce() bool {
	for _, card := range p.Hand {
		if card.Value == "Ace" {
			return true
		}
	}
	return false
}

// Busted evaluates whether or not a player has busted.
func (p *Player) Busted() bool {
	return p.Score > 21
}
