package main

import (
	"fmt"

	"github.com/fuskovic/gophercises/deck"
)

type (
	player struct {
		name                string
		score, wins, losses uint
		//bet uint
		hand []*deck.Card
	}
	players []*player
)

func (p *player) hit(newCard *deck.Card) {
	fmt.Printf("%s hit a %s of %s\n", p.name, newCard.Value, newCard.Suit)

	p.hand = append(p.hand, newCard)

	handleAce := func() {
		if p.score < 22 && (p.score+11) < 22 {
			p.score += 11
		} else if p.score < 22 && (p.score+11) > 21 {
			p.score += 1
		}
	}

	if newCard.Value == "Ace" {
		handleAce()
	} else {
		p.score += uint(deck.CardMap[newCard.Value])
	}

	fmt.Printf("%s's score : %d\n", p.name, p.score)
}

func (p *player) getAutomatedInput() int {
	if p.score <= 16 || p.score == 17 && p.hasAce() {
		return hit
	}
	return stand
}

func (p *player) hasAce() bool {
	for _, card := range p.hand {
		if card.Value == "Ace" {
			return true
		}
	}
	return false
}
