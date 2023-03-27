package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/movies", handleMovies)
	http.ListenAndServe("localhost:3000", r)
}

func handleMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get movies")
}
