package creating_test

import (
	"testing"

	"github.com/bennyc/deck/internal/creating"
	"github.com/bennyc/deck/internal/memory"
)

func TestUnitServiceCreatesNewDeck(t *testing.T) {
	s := creating.New(memory.New())
	deck, err := s.New(creating.Options{})

	if err != nil {
		t.Errorf("service could not be created with default options")
	}

	if len(deck.Cards) != 52 {
		t.Errorf("count of cards in deck expected to be: %d, got %d", 52, len(deck.Cards))
	}

	if deck.Id == "" {
		t.Errorf("deck.Id has not been set after being created")
	}
}

func TestUnitServiceValidationForBadCard(t *testing.T) {
	s := creating.New(memory.New())
	_, err := s.New(creating.Options{
		Selection: []string{"INCORRECTCARD"},
	})

	if err == nil {
		t.Errorf("deck creation was allowed with invalid card code")
	}
}
