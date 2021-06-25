package drawing_test

import (
	"testing"

	"github.com/bennyc/deck/internal/drawing"
	"github.com/bennyc/deck/internal/entity"
	"github.com/bennyc/deck/internal/memory"
)

func TestUnitCannotDrawMoreCardsThanInDeck(t *testing.T) {
	mem := memory.New()
	deck := entity.NewStandardDeck()
	mem.Save(deck)

	s := drawing.New(mem)
	if _, err := s.Draw(deck, len(deck.Cards)+1); err == nil {
		t.Errorf("service allows consumer to draw more cards than available in deck")
	}
}
