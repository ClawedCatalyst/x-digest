package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	httpadapter "xdigest/internal/adapter/inbound/http"
	"xdigest/internal/adapter/outbound/crypto"
	"xdigest/internal/adapter/outbound/postgres"
	"xdigest/internal/adapter/outbound/xapi"
	"xdigest/internal/application"
	"xdigest/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	crypter, err := crypto.NewCrypter(cfg.EncKey)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepo(db.Pool)
	tokenRepo := postgres.NewTokenRepo(db.Pool)
	digestRepo := postgres.NewDigestRepo(db.Pool)

	xClient := xapi.NewClient(xapi.Config{
		XClientID:     cfg.XClientID,
		XClientSecret: cfg.XClientSecret,
		XRedirectURI:  cfg.XRedirectURI,
	})

	authSvc := application.NewAuthService(application.AuthServiceConfig{
		XClientID:       cfg.XClientID,
		XRedirectURI:    cfg.XRedirectURI,
		XScopes:         cfg.XScopes,
		FrontendBaseURL: cfg.FrontendBaseURL,
	}, crypter, userRepo, tokenRepo, xClient)

	digestSvc := application.NewDigestService(digestRepo, tokenRepo, crypter, xClient)

	server := httpadapter.NewServer(authSvc, digestSvc, cfg.FrontendBaseURL)

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.FrontendBaseURL},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	r.Get("/auth/x/start", server.HandleAuthStart)
	r.Get("/auth/x/callback", server.HandleAuthCallback)
	r.Get("/digest", server.HandleGetDigestPeriod)
	r.Get("/digest/today", server.HandleGetDigestToday)
	r.Post("/jobs/digest", server.HandleBuildDigestPeriod)
	r.Post("/jobs/digest/today", server.HandleBuildDigestToday)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("backend listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
