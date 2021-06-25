package memory

import (
	"errors"

	"github.com/bennyc/deck/internal/entity"
	"github.com/google/uuid"
)

type storage struct {
	decks map[string]*entity.Deck
}

// Save Deck will associate an ID to the Deck and then place within the
// Repository of Decks!
func (s *storage) Save(d *entity.Deck) error {
	d.Id = uuid.NewString()
	s.decks[d.Id] = d

	return nil
}

// GetById returns a pointer to a Deck based on the ID provided
func (s *storage) GetById(id string) (*entity.Deck, error) {
	if val, ok := s.decks[id]; ok {
		return val, nil
	}

	return nil, errors.New("could not find deck")
}

// Update will do nothing in the context of memory storage
// As we're using a pointer to a Deck, any updates will be immediately reflected
// within the Repository and don't need to be committed!
func (s *storage) Update(*entity.Deck) error {
	return nil
}

func New() entity.DeckRepository {
	return &storage{
		decks: make(map[string]*entity.Deck),
	}
}
