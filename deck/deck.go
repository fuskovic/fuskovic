package deck

var (
	// Suits contains all suits that a playing card could have.
	Suits = []string{"spades", "diamonds", "clubs", "hearts"}

	CardMap = map[string]int{
		"Joker": 0,
		"Ace":   1,
		"Two":   2,
		"Three": 3,
		"Four":  4,
		"Five":  5,
		"Six":   6,
		"Seven": 7,
		"Eight": 8,
		"Nine":  9,
		"Ten":   10,
		"Jack":  10,
		"Queen": 10,
		"King":  10,
	}
)

type (
	// Card contains a suit and rank representing a traditional playing card.
	Card struct{ Value, Suit string }

	// Deck implements the sort pkgs Interface interface for a slice of playing cards.
	Deck []Card
)

// New initializes a new deck.
func New() *Deck {
	var deck Deck
	for _, suit := range Suits {
		for name := range CardMap {
			if name == "Joker" {
				continue
			}
			deck = append(deck, Card{
				Value: name,
				Suit:  suit,
			})
		}
	}
	deck = Shuffle(deck)
	return &deck
}

func (d *Deck) Len() int {
	return len(*d)
}

func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Draw removes a card from the deck, returning the removed card and the deck with the card removed.
func (d *Deck) Draw() (Card, *Deck) {
	deckCopy := *d
	lastCard := deckCopy[d.Len()-1]
	updatedDeck := d.RemoveCard(lastCard)
	return lastCard, &updatedDeck
}

// RemoveCard removes a card from the deck.
func (d *Deck) RemoveCard(c Card) (new Deck) {
	for _, card := range *d {
		if card == c {
			continue
		}
		new = append(new, card)
	}
	return
}

// GetCardIndez returns the position of a particular card in the deck.
func (d *Deck) GetCardIndex(c Card) (cardIndex int) {
	for i, card := range *d {
		if card == c {
			cardIndex = i
		}
	}
	return
}
