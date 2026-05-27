package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	db "github.com/bradonstudie/lepidoptera/internal/db/generated"
	"github.com/bradonstudie/lepidoptera/internal/handlers"
	"github.com/bradonstudie/lepidoptera/internal/mailer"
	"github.com/bradonstudie/lepidoptera/internal/service"
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

	// services
	showService := service.NewShowService(queries)
	subscriberService := service.NewSubscriberService(queries, m, tokenSecret)
	adminService := service.NewAdminService(queries, m)

	// worker
	worker := mailer.NewWorker(m, queries, tokenSecret)

	// middleware
	adminMiddleware := handlers.NewAdminMiddleware(queries)

	s := &Server{
		router: chi.NewRouter(),
		Worker: worker,
	}

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(httprate.Limit(
		10,
		10*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	))

	// Login rate limiters
	// SEE: https://github.com/go-chi/httprate
	loginRateLimit := httprate.Limit(5, time.Minute, httprate.WithKeyFuncs(httprate.KeyByIP))
	loginEmailRateLimit := httprate.Limit(3, 10*time.Minute, httprate.WithKeyFuncs(
		func(r *http.Request) (string, error) {
			r.ParseForm()
			if email := r.FormValue("email"); email != "" {
				return "login:email:" + email, nil
			}
			return httprate.KeyByIP(r)
		},
	))

	s.router.Handle("/static/*", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

	// handlers
	shows := handlers.NewShowHandler(showService)
	subs := handlers.NewSubscriberHandler(subscriberService)
	admin := handlers.NewAdminHandler(adminService, queries, tokenSecret)

	// public routes
	s.router.Get("/", shows.Index)
	s.router.Get("/shows/{slug}", shows.Detail)
	s.router.Get("/subscribe", subs.SubscribePage)
	s.router.Post("/subscribe", subs.Subscribe)
	s.router.Get("/confirm", subs.Confirm)
	s.router.Get("/unsubscribe", subs.Unsubscribe)

	// admin auth — public
	s.router.Get("/admin/login", admin.LoginPage)
	s.router.With(loginRateLimit, loginEmailRateLimit).Post("/admin/login", admin.Login)
	s.router.Get("/admin/verify", admin.Verify)

	// admin routes — protected
	s.router.Group(func(r chi.Router) {
		r.Use(adminMiddleware.AdminOnly)
		r.Post("/admin/logout", admin.Logout)
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
