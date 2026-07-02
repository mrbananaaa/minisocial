package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	postapi "github.com/mrbananaaa/minisocial/internal/post/api"
	userapi "github.com/mrbananaaa/minisocial/internal/user/api"
)

type Routes struct {
	userHandler *userapi.Handler
	postHandler *postapi.Handler
}

func NewRouter(h Routes) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromHeader("X-Real-IP"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// TODO: change cors origin
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Compress(5, "application/json"))

	r.Route("/users", func(u chi.Router) {
		u.Post("/", h.userHandler.CreateUser)
		u.Get("/{id}", h.userHandler.GetUser)
		u.Patch("/{id}", h.userHandler.UpdateProfile)
		u.Delete("/{id}", h.userHandler.DeleteUser)
	})

	r.Route("/posts", func(u chi.Router) {
		u.Post("/", h.postHandler.CreatePost)
		u.Patch("/{id}", h.postHandler.EditPost)
		u.Delete("/{id}", h.postHandler.ArchivePost)
	})

	return r
}
