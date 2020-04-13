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
	Deck []*Card
)

// New initializes a new deck.
func New() *Deck {
	var deck Deck
	for _, suit := range Suits {
		for name := range CardMap {
			if name == "Joker" {
				continue
			}
			deck = append(deck, &Card{
				Value: name,
				Suit:  suit,
			})
		}
	}
	return Shuffle(&deck)
}

// Len returns the number of cards in the deck.
func (d *Deck) Len() int {
	return len(*d)
}

// Swap replaces the positions of two cards by the indices provided.
func (d *Deck) Swap(i, j int) {
	deck := *d
	deck[i], deck[j] = deck[j], deck[i]
	*d = deck
}

// Draw removes the last card from the deck and returns it.
func (d *Deck) Draw() (*Card) {
	lastCard := d.GetCard(d.Len()-1)
	d.RemoveCard(lastCard)
	return lastCard
}

// RemoveCard removes a card from the deck.
func (d *Deck) RemoveCard(c *Card) {
	var new Deck	
	for _, card := range *d {
		if card != c {
			new = append(new, card)	
		}
	}
	*d = new
}

// GetCard returns a card from the deck by its index without removing it from the deck.
func (d *Deck) GetCard(index int) *Card {
	for i, c := range *d {
		if i == index {
			return c
		}
	}
	return nil
}
