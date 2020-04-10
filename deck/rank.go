package deck

// Represents the numerical value of a playing card.
const (
	Face Rank = Ten

	Joker Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack  = Face
	Queen = Face
	King  = Face
)

// Rank represents the numerical value of a playing card.
type Rank uint

// Rankings returns all Ranks
func Rankings() (rankings []Rank) {
	for i := 2; i < 15; i++ {
		rankings = append(rankings, Rank(i))
	}
	return
}
