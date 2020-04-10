package deck

import (
	"math/rand"
	"time"
)

// TODO :
// - Create a type-alias for a []Card (Deck?) that implements the
// Interface interface in the sort pkg https://godoc.org/sort#Interface
// - Hang two methods off the new type to satisfy
//   - An option to sort the cards with a user-defined comparison function.
//		The sort package in the standard library can be used here, and expects a less(i, j int) bool function.
//	 - A default comparison function that can be used with the sorting option.

// SortFunc describes the method signature for a functional option that can
// be passed to Sort.
type SortFunc func(cards []Card) []Card

// Shuffle rearranges a deck of cards in a random order.
var Shuffle SortFunc = func(cards []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i int, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// Sort applies a given set of sort functions to a deck of playing cards.
func Sort(cards []Card, funcs ...SortFunc) []Card {
	for _, f := range funcs {
		cards = f(cards)
	}
	return cards
}

// AddJokers returns a SortFunc used to add n number of jokers to a deck.
func AddJokers(numJokersToAdd int, cards []Card) SortFunc {
	return func(cards []Card) []Card {
		for i := 1; i < numJokersToAdd; i++ {
			cards = append(cards, Card{
				Rank: Joker,
				Suit: "Joker",
			})
		}
		return cards
	}
}

// FilterByRank returns a SortFunc that removes all cards from a deck for all given ranks.
func FilterByRank(cards []Card, ranks ...Rank) SortFunc {
	return func(cards []Card) []Card {
		var filtered []Card
		counts := make(map[Card]int)

		for _, r := range ranks {
			for _, card := range cards {
				counts[card]++

				if card.Rank == r {
					continue
				}

				if counts[card] == 1 {
					filtered = append(filtered, card)
				}
			}
		}
		return filtered
	}
}

// Multi returns a SortFunc that adds n number of decks to an existing deck
func Multi(cards []Card, numOfDecks int) SortFunc {
	return func(cards []Card) []Card {
		var deck []Card
		decks := [][]Card{cards}

		for i := 1; i < numOfDecks; i++ {
			decks = append(decks, Shuffle(cards))
		}

		for _, d := range decks {
			for _, card := range d {
				deck = append(deck, card)
			}
		}
		return deck
	}
}
