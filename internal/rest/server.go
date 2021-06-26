package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/bennyc/deck/internal/creating"
	"github.com/bennyc/deck/internal/drawing"
	"github.com/bennyc/deck/internal/entity"
	"github.com/bennyc/deck/internal/err"
	"github.com/gorilla/mux"
)

type Server struct {
	Creating   creating.Service
	Drawing    drawing.Service
	Repository entity.DeckRepository
}

func (s *Server) Routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/decks", createDeck(s.Creating)).Methods("POST")
	r.HandleFunc("/decks/{id}", openDeck(s.Repository)).Methods("GET")
	r.HandleFunc("/decks/{id}/draw-cards", drawCards(s.Drawing, s.Repository)).Methods("POST")

	r.Use(contentType)
	r.Use(logRequest)

	return r
}

func contentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := struct {
			start  time.Time
			method string
			uri    string
		}{
			time.Now(),
			r.Method,
			r.RequestURI,
		}

		log.Printf("[%s] %s %s", fields.start, fields.method, fields.uri)
		next.ServeHTTP(w, r)
	})
}

func renderJSON(w http.ResponseWriter, r interface{}, s int) {
	str, err := json.Marshal(&r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(s)
	w.Write(str)
}

func renderError(w http.ResponseWriter, e error) {
	if statusErr := err.ErrStatusCode(nil); errors.As(e, &statusErr) {
		log.Println(statusErr)

		code := statusErr.Code()
		w.WriteHeader(code)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
