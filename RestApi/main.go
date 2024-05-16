package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})

		r.Get("/article", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("this is the article"))
		})

		r.Route("/hello", func(r chi.Router) {
			r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
				name := chi.URLParam(r, "name")
				w.Write([]byte("hello " + name))
			})
		})
	})

	http.ListenAndServe(":3000", r)
}
