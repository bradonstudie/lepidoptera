package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	db "github.com/lepidoptera/lepidoptera/internal/db/generated"
	"github.com/lepidoptera/lepidoptera/internal/handlers"
	"github.com/lepidoptera/lepidoptera/internal/mailer"
)

type Server struct {
	router *chi.Mux
	Worker *mailer.Worker
}

func New(dbURL string, m mailer.Mailer, tokenSecret string) (*Server, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("db connect: %w", err)
	}

	sqlDB := stdlib.OpenDBFromPool(pool)
	queries := db.New(sqlDB)
	worker := mailer.NewWorker(m, queries, tokenSecret)

	s := &Server{
		router: chi.NewRouter(),
		Worker: worker,
	}

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	// static files
	s.router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

	// handlers
	shows := handlers.NewShowHandler(queries)
	subs := handlers.NewSubscriberHandler(queries, m, tokenSecret)
	admin := handlers.NewAdminHandler(queries, m, tokenSecret)

	// middleware
	adminMiddleware := handlers.NewAdminMiddleware(queries)

	// public routes
	s.router.Get("/", shows.Index)
	s.router.Get("/shows/{slug}", shows.Detail)
	s.router.Get("/subscribe", subs.SubscribePage)
	s.router.Post("/subscribe", subs.Subscribe)
	s.router.Get("/confirm", subs.Confirm)
	s.router.Get("/unsubscribe", subs.Unsubscribe)

	s.router.Get("/admin/login", admin.LoginPage)
	s.router.Post("/admin/login", admin.Login)
	s.router.Post("/admin/logout", admin.Logout)
	s.router.Get("/admin/verify", admin.Verify)

	// admin routes
	s.router.Group(func(r chi.Router) {
		r.Use(adminMiddleware.AdminOnly)
		r.Get("/admin", admin.Dashboard)
		r.Get("/admin/shows/new", admin.NewShowForm)
		r.Post("/admin/shows", admin.CreateShow)
		r.Post("/admin/shows/{id}/publish", admin.PublishShow)
		r.Get("/admin/bands/new", admin.NewBandForm)
		r.Post("/admin/bands", admin.CreateBand)
		r.Get("/admin/venues/new", admin.NewVenueForm)
		r.Post("/admin/venues", admin.CreateVenue)
	})

	return s, nil
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
