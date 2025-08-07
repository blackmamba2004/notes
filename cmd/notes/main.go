package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"notes/internal/config"
	"notes/internal/lib/logger/sl"

	"notes/internal/notes/handlers"
	"notes/internal/storage/postgresql"
	"os"
	"time"
)

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			entry := logger.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			rw := &customeResponseWriter{
				ResponseWriter: w,
			}

			start := time.Now()
			defer func() {
				duration := time.Since(start)
				entry.Info("request completed",
					slog.Int("status", rw.statusCode),
					slog.String("duration", duration.String()),
				)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}


const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	logger.Info("notes-logger started", slog.String("env", cfg.Env))

	storage, err := postgresql.New(cfg.Database.BuildPostgresDSN())

	if err != nil {
		logger.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	defer storage.Close()

	mux := http.NewServeMux()

	middlewareHandler := loggingMiddleware(logger)(mux)

	server := http.Server{
		Handler: middlewareHandler,
		Addr:    "localhost:8083",
	}

	fmt.Printf("Server started on http://%s\n", server.Addr)

	mux.HandleFunc("/", Hello)
	mux.HandleFunc("/notes", handlers.CreateNote)

	server.ListenAndServe()
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"error": "not found",
	})
}

type customeResponseWriter struct {
	http.ResponseWriter
	statusCode int
	set    bool
}

func (rw *customeResponseWriter) WriteHeader(code int) {
	if !rw.set {
		rw.statusCode = code
		rw.set = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *customeResponseWriter) Write(b []byte) (int, error) {
	if !rw.set {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(
				os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}
	return logger
}
