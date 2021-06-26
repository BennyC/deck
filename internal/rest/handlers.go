package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/bennyc/deck/internal/creating"
	"github.com/bennyc/deck/internal/drawing"
	"github.com/bennyc/deck/internal/entity"
	derr "github.com/bennyc/deck/internal/err"
	"github.com/gorilla/mux"
)

// Create Deck handler for REST implementation
// Will use creating service to generate the deck and persist it to the desired
// storage
func createDeck(service creating.Service) http.HandlerFunc {
	type createReq struct {
		Shuffle   bool   `json:"shuffle"`
		Selection string `json:"cards"`
	}

	type createRes struct {
		Id        string `json:"deck_id"`
		Shuffled  bool   `json:"shuffled"`
		Remaining int    `json:"remaining"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var cr createReq

		if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
			renderError(w, derr.New(http.StatusBadRequest, err))
			return
		}

		// Split the string on commas, and delegate the selection to our
		// creating Service, which will perform validation on the values provided
		sel := []string{}
		if len(cr.Selection) > 0 {
			sel = strings.Split(cr.Selection, ",")
		}

		deck, err := service.New(creating.Options{
			Shuffle:   cr.Shuffle,
			Selection: sel,
		})

		if err != nil {
			if errors.Is(err, creating.ErrOptionValidation) {
				err = derr.New(http.StatusUnprocessableEntity, err)
			}

			renderError(w, err)
			return
		}

		w.Header().Set("Location", r.URL.Path+"/"+deck.Id)
		renderJSON(w, createRes{
			Id:        deck.Id,
			Shuffled:  deck.Shuffled,
			Remaining: len(deck.Cards),
		}, http.StatusCreated)
	}
}

// Attempt to open a requested deck
func openDeck(repository entity.DeckRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deck, err := repository.GetById(vars["id"])

		if err != nil {
			renderError(w, derr.New(http.StatusNotFound, err))
			return
		}

		renderJSON(w, deck, http.StatusOK)
	}
}

func drawCards(service drawing.Service, repository entity.DeckRepository) http.HandlerFunc {
	type request struct {
		Count int `json:"count"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		deck, err := repository.GetById(vars["id"])

		if err != nil {
			renderError(w, derr.New(http.StatusNotFound, err))
			return
		}

		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			renderError(w, derr.New(http.StatusBadRequest, err))
			return
		}

		drawn, err := service.Draw(deck, req.Count)
		if err != nil {
			renderError(w, derr.New(http.StatusUnprocessableEntity, err))
			return
		}

		renderJSON(w, struct {
			Cards []entity.Card `json:"cards"`
		}{
			drawn,
		}, http.StatusOK)
	}
}
