package creating

import (
	"github.com/bennyc/deck/internal/entity"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Create a new Creating Service
func New(decks entity.DeckRepository) Service {
	return &service{
		decks: decks,
	}
}

type Service interface {
	// Create a new Deck based on Options provided and will be stored in the DeckRepository
	// Options.Selection allows the consumer to declare a specific list of cards they'd like to use
	// Options.Shuffle allows the consumer to shuffle the deck prior to saving
	New(Options) (*entity.Deck, error)
}

type service struct {
	decks entity.DeckRepository
}

// Create a new Deck
func (s *service) New(o Options) (*entity.Deck, error) {
	// Validate the Options object to ensure all Selection Codes exist
	// within a Standard Deck
	if err := o.Validate(); err != nil {
		return nil, err
	}

	deck := entity.NewStandardDeck()

	// If there is a Selection of Cards, filter out all of the cards
	// that do not exist in the Selection and place the Selected cards
	// back into the Deck
	if len(o.Selection) > 0 {
		cards := []entity.Card{}
		for _, c := range deck.Cards {
			if Contains(o.Selection, c.Code) {
				cards = append(cards, c)
			}
		}

		deck.Cards = cards
	}

	// If required, we shuffle the deck before persist the Deck
	if o.Shuffle {
		deck.Shuffle()
	}

	if err := s.decks.Save(deck); err != nil {
		return nil, err
	}

	return deck, nil
}

type Options struct {
	// Should the Deck be shuffled after creation
	Shuffle bool

	// Choice of card codes, all other cards will not be used in the deck
	// that is created
	Selection []string
}

// Validate the Selection option is in the allowed card codes
// https://github.com/go-ozzo/ozzo-validation/issues/82
func (o Options) Validate() error {
	deck := entity.NewStandardDeck()
	allowed := []interface{}{}
	for _, card := range deck.Cards {
		allowed = append(allowed, card.Code)
	}

	return validation.ValidateStruct(
		&o,
		validation.Field(&o.Selection, validation.Each(validation.In(allowed...))),
	)
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
