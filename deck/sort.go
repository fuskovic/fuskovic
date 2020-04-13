package deck

import (
	"math/rand"
	"time"
)

// SortFunc describes the method signature for a functional option that can be passed to Sort.
type SortFunc func(deck *Deck) *Deck

// Shuffle rearranges a deck of cards in a random order.
var Shuffle SortFunc = func(d *Deck) *Deck {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(d.Len(), func(i int, j int) { d.Swap(i, j) })
	return d
}

// Sort applies a given set of sort functions to a deck of playing cards.
func Sort(d *Deck, funcs ...SortFunc) *Deck {
	for _, f := range funcs {
		d = f(d)
	}
	return d
}

// AddJokers returns a SortFunc used to add n number of jokers to a deck.
func AddJokers(numJokersToAdd int, cards Deck) SortFunc {
	return func(d *Deck) *Deck {
		var deck Deck
		for i := 1; i < numJokersToAdd; i++ {
			deck = append(deck, &Card{
				Value: "Joker",
				Suit:  "Joker",
			})
		}
		return &deck
	}
}

// FilterByRank returns a SortFunc that removes all cards from a deck for all given ranks.
func FilterByRank(d Deck, ranks ...string) SortFunc {
	return func(d *Deck) *Deck {
		var filtered Deck
		counts := make(map[Card]int)

		for _, r := range ranks {
			for _, card := range *d {
				counts[*card]++

				if counts[*card] == 1 && card.Value != r {
					filtered = append(filtered, card)
				}
			}
		}
		return &filtered
	}
}

// Multi returns a SortFunc that adds n number of decks to an existing deck
func Multi(d *Deck, numOfDecks int) SortFunc {
	return func(d *Deck) *Deck {
		var deck Deck
		deck = append(deck, *Shuffle(d)...)
		return &deck
	}
}