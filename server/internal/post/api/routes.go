package api

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", h.CreatePost)
		r.Patch("/{id}", h.EditPost)
		r.Delete("/{id}", h.ArchivePost)
	})
}
