package api

import "github.com/go-chi/chi/v5"

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
		r.Patch("/{id}", h.UpdateProfile)
		r.Delete("/{id}", h.DeleteUser)
	})
}
