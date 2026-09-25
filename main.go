package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/sahilverma/muze-go-backend/auth"
	"github.com/sahilverma/muze-go-backend/discussions"
	"github.com/sahilverma/muze-go-backend/platform"
)

//go:embed web/*
var webFiles embed.FS

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	store := discussions.NewMemoryStore()
	service := discussions.NewService(store, platform.NewNoopCache(), platform.NewLogPublisher(logger))
	tokens := auth.NewTokenManager(env("MUZE_JWT_SECRET", "local-development-secret"), 24*time.Hour)

	api := platform.RequestID(platform.Logging(logger, discussions.NewHandler(service, tokens)))
	web, err := fs.Sub(webFiles, "web")
	if err != nil {
		logger.Error("load_web", "error", err)
		os.Exit(1)
	}
	static := http.FileServer(http.FS(web))
	root := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || len(r.URL.Path) >= 8 && r.URL.Path[:8] == "/api/v1/" {
			api.ServeHTTP(w, r)
			return
		}
		static.ServeHTTP(w, r)
	})

	addr := env("PORT", "8080")
	logger.Info("server_started", "address", "http://localhost:"+addr)
	if err := http.ListenAndServe(":"+addr, root); err != nil {
		logger.Error("server_stopped", "error", err)
		os.Exit(1)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}