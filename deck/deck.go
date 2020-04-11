package deck

// Suits contains all suits that a playing card could have.
var Suits = []string{"spades", "diamonds", "clubs", "hearts"}

type (
	// Card contains a suit and rank representing a traditional playing card.
	Card struct {
		Rank Rank
		Suit string
	}
	// Deck implements the sort pkgs Interface interface for a slice of playing cards.
	Deck []Card
)

// New initializes a new deck.
func New() Deck {
	var deck Deck
	for _, suit := range Suits {
		for _, rank := range Rankings() {
			deck = append(deck, Card{
				Rank: rank,
				Suit: suit,
			})
		}
	}
	return Shuffle(deck)
}
