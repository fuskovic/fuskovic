package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/fuskovic/gophercises/deck"
)

const (
	waiting = iota
	stand
	hit
)

// Game represents a game of Blackjack.
type Game struct {
	deck    *deck.Deck
	players *players
	winners *players
}

// NewGame starts a new game of Blackjack by adding n players to the game.
// New Game also gives each player 1 card once per round for two rounds.
// If anyone hits 21 after those two rounds, then win automatically.
func NewGame(n int) *Game {
	var dealer = &player{name: "dealer"}
	var players, winners players

	for i := 1; i < n; i++ {
		newPlayer := player{name: fmt.Sprintf("player%d", i)}
		players = append(players, &newPlayer)
	}

	players = append(players, dealer)
	d := deck.New()

	// give each player two cards
	for i := 0; i < 2; i++ {
		for _, p := range players {
			p.hit(d.Draw())
			if p.score == 21 {
				winners = append(winners, p)
			}
		}
	}

	return &Game{
		deck:    d,
		players: &players,
		winners: &winners,
	}
}

// Play runs the game, returning all winners.
func (g *Game) Play() players {
	if len(*g.winners) > 0 {
		return *g.winners
	}

	selection := waiting

game:
	for _, p := range *g.players {
		for selection == waiting || selection == hit {
			selection = g.getSelection(p)

			if selection == hit {
				fmt.Printf("%s has chose to hit\n", p.name)
				p.hit(g.deck.Draw())

				if p.score > 21 {
					g.handleLoser(p.name)
					if p.name == "player1" || p.name == "dealer" {
						break game
					}
					continue
				}

				if p.score == 21 {
					g.incrementWins(p.name)
					break game
				}
			} else {
				fmt.Printf("%s has chose to stay\n", p.name)
			}
		}
		selection = waiting
	}
	return g.getWinners()
}

func (g *Game) getSelection(p *player) int {
	if p.name == "player1" {
		fmt.Println("Your turn\n1 - stay\n2 - hit")
		return g.getUserInput()
	}

	fmt.Printf("its %s's turn\n", p.name)
	return p.getAutomatedInput()
}

func (g *Game) getUserInput() (selection int) {
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

// handleLoser increments the number of losses for a particular player.
//
// If the busted player == dealer then player1's wins are incremented.
// If the busted player == player1 then the dealers wins are incremented.
func (g *Game) handleLoser(name string) {
	var winner, loser string

	switch name {
	case "dealer":
		fmt.Println("dealer has busted - YOU WIN!")
		winner, loser = "player1", "dealer"
	case "player1":
		fmt.Println("you busted - dealer wins")
		winner, loser = "dealer", "player1"
	default:
		fmt.Printf("%s has busted\n%s is out!", name, name)
		g.incrementLosses(name)
	}

	if winner != "" {
		g.incrementWins(winner)
	}

	if loser != "" {
		g.incrementLosses(loser)
	}
}

func (g *Game) incrementLosses(name string) {
	for _, p := range *g.players {
		if p.name == name {
			p.losses++
		}
	}
}

func (g *Game) incrementWins(name string) {
	var winners players

	for _, p := range *g.players {
		if p.name == name {
			p.wins++
			winners = append(winners, p)
		}
	}
	g.winners = &winners
}

// Reset returns all players hands to the deck, resets their scores and resets the winner.
func (g *Game) Reset() {
	var winners players
	fmt.Println("********* RESETTING....**********")

	for _, p := range *g.players {
		p.score = 0
		p.hand = []*deck.Card{}
	}
	g.deck = deck.New()

	// give each player two cards
	for i := 0; i < 2; i++ {
		for _, p := range *g.players {
			p.hit(g.deck.Draw())
			if p.score == 21 {
				winners = append(winners, p)
			}
		}
	}
	g.winners = &winners
}

func (g *Game) getWinners() (winners players) {
	if len(*g.winners) > 0 {
		return *g.winners
	}

	var highestScore uint

	scores := make(map[uint][]string)

	for _, p := range *g.players {
		scores[p.score] = append(scores[p.score], p.name)
	}

	for score := range scores {
		if score > highestScore {
			highestScore = score
		}
	}

	for _, winner := range scores[highestScore] {
		winners = append(winners, &player{name: winner})
	}
	return
}
