package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luismingati/buymeacoffee/internal/config"
	"github.com/luismingati/buymeacoffee/internal/service"
	"github.com/luismingati/buymeacoffee/internal/store/pg"
)

type apiConfig struct {
	q     *pg.Queries
	r     *chi.Mux
	pool  *pgxpool.Pool
	jwt   *service.JWTService
	m     service.Mailer
	redis *service.RedisService
}

func (cfg *apiConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cfg.r.ServeHTTP(w, r)
}

func ApiHandler(q *pg.Queries, pool *pgxpool.Pool) http.Handler {
	ctx := context.Background()
	jwt := service.NewJWTService(config.GetSecretKey())

	m, err := service.NewSESEmailer(ctx, "contato@luismingati.dev")
	if err != nil {
		panic(err)
	}

	redis, err := service.NewRedisService(ctx)
	if err != nil {
		panic(err)
	}

	cfg := &apiConfig{
		q:     q,
		pool:  pool,
		jwt:   jwt,
		m:     m,
		redis: redis,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthcheck", cfg.HealthcheckHandler)
	r.Route("/api", func(r chi.Router) {
		r.Post("/signup", cfg.SignupHandler)
		r.Post("/signin", cfg.SigninHandler)
		r.Post("/forgot-password", cfg.ForgotPasswordHandler)
		r.Post("/reset-password", cfg.ResetPasswordHandler)
	})

	cfg.r = r
	return cfg
}

func (cfg *apiConfig) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
