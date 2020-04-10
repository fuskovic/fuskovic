package deck

// Suits contains all suits that a playing card could have.
var Suits = []string{"spades", "diamonds", "clubs", "hearts"}

// Card contains a suit and rank representing a traditional playing card.
type Card struct {
	Rank Rank
	Suit string
}

// New initializes a new deck.
func New() []Card {
	var deck []Card
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
