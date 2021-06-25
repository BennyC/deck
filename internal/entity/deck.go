package entity

import (
	"encoding/json"
	"math/rand"
	"time"
)

type DeckRepository interface {
	Save(*Deck) error
	Update(*Deck) error
	GetById(string) (*Deck, error)
}

type Card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
	Code  string `json:"code"`
}

type Deck struct {
	Id       string
	Cards    []Card
	Shuffled bool
}

// MarshalJSON will include a count of the remaining cards, which is computed
// via the length from the Deck
func (d *Deck) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id        string `json:"id"`
		Cards     []Card `json:"cards"`
		Shuffled  bool   `json:"shuffled"`
		Remaining int    `json:"remaining"`
	}{
		d.Id,
		d.Cards,
		d.Shuffled,
		len(d.Cards),
	})
}

// Shuffle the deck of cards and flag that the deck has been shuffled
// Will use new seed per shuffle
func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i int, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})

	d.Shuffled = true
}

// NewStandardDeck
func NewStandardDeck() *Deck {
	cards := []Card{}
	suits := []string{"CLUBS", "DIAMONDS", "HEARTS", "SPADES"}
	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	for _, suit := range suits {
		for _, rank := range ranks {
			cards = append(cards, Card{
				Suit:  suit,
				Value: rank,
				Code:  rank + string(suit[0]),
			})
		}
	}

	return &Deck{Cards: cards}
}
