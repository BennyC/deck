package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bennyc/deck/internal/creating"
	"github.com/bennyc/deck/internal/drawing"
	"github.com/bennyc/deck/internal/memory"
	"github.com/bennyc/deck/internal/rest"
)

func main() {
	mem := memory.New()
	s := &rest.Server{
		Creating:   creating.New(mem),
		Drawing:    drawing.New(mem),
		Repository: mem,
	}

	var port string
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Println("Starting server on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Routes()))
}
