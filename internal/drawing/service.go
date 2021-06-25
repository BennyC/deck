package drawing

import (
	"errors"

	"github.com/bennyc/deck/internal/entity"
)

type Service interface {
	Draw(*entity.Deck, int) ([]entity.Card, error)
}

type service struct {
	r entity.DeckRepository
}

// Create a new Drawing Service
func New(r entity.DeckRepository) Service {
	return &service{r}
}

// Draw cards from a Deck
// If the Deck does not contain enough cards an error will be returned
// If the Deck contains enough Cards, we return the count of Cards and remove
// those from the Deck
func (s service) Draw(d *entity.Deck, count int) ([]entity.Card, error) {
	if len(d.Cards) < count {
		return nil, errors.New("not enough cards to draw")
	}

	drawn := d.Cards[:count]
	d.Cards = d.Cards[count:]
	return drawn, nil
}
