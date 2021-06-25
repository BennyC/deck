package rest_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bennyc/deck/internal/creating"
	"github.com/bennyc/deck/internal/drawing"
	"github.com/bennyc/deck/internal/entity"
	"github.com/bennyc/deck/internal/memory"
	"github.com/bennyc/deck/internal/rest"
	"github.com/gavv/httpexpect/v2"
)

func createServer() (*rest.Server, entity.DeckRepository) {
	mem := memory.New()
	s := &rest.Server{
		Creating:   creating.New(mem),
		Drawing:    drawing.New(mem),
		Repository: mem,
	}

	return s, mem
}

func TestIntegrationSuccessfulCreateDeck(t *testing.T) {
	testCases := []struct {
		Shuffle  bool   `json:"shuffle"`
		Cards    string `json:"cards"`
		deckSize int
	}{
		{
			Shuffle:  true,
			Cards:    "10S,AD",
			deckSize: 2,
		},

		{
			Shuffle:  false,
			Cards:    "",
			deckSize: 52,
		},
	}

	srv, _ := createServer()
	server := httptest.NewServer(srv.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	for _, tc := range testCases {
		res := e.POST("/decks").WithJSON(tc)
		exp := res.Expect()

		exp.Header("Location")
		exp.Status(http.StatusCreated)

		obj := exp.JSON().Object()
		obj.Keys().Contains("deck_id", "shuffled", "remaining")
		obj.ContainsKey("shuffled").ValueEqual("shuffled", tc.Shuffle)
		obj.ContainsKey("remaining").ValueEqual("remaining", tc.deckSize)
	}
}

func TestIntegrationViewDeck(t *testing.T) {
	srv, mem := createServer()
	server := httptest.NewServer(srv.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)

	deck := entity.NewStandardDeck()
	mem.Save(deck)

	res := e.GET("/decks/" + deck.Id)
	exp := res.Expect()
	exp.Status(http.StatusOK)

	obj := exp.JSON().Object()
	obj.Keys().Contains("id", "remaining", "shuffled", "cards")
	obj.ContainsKey("shuffled").ValueEqual("shuffled", deck.Shuffled)
	obj.ContainsKey("remaining").ValueEqual("remaining", len(deck.Cards))
}

func TestIntegrationDrawFromDeck(t *testing.T) {
	srv, mem := createServer()
	server := httptest.NewServer(srv.Routes())
	defer server.Close()

	e := httpexpect.New(t, server.URL)
	deck := entity.NewStandardDeck()
	mem.Save(deck)

	willDraw := deck.Cards[0:3]
	log.Println(willDraw)

	req := map[string]interface{}{"count": 3}
	res := e.POST("/decks/" + deck.Id + "/draw-cards").WithJSON(req)

	exp := res.Expect()
	exp.Status(http.StatusOK)

	obj := exp.JSON().Object()
	obj.Keys().Contains("cards")
	obj.Value("cards").Array().Length().Equal(req["count"])

	for i, c := range willDraw {
		obj.Value("cards").
			Array().
			Element(i).
			Object().
			ContainsKey("code").
			ValueEqual("code", c.Code)
	}
}
