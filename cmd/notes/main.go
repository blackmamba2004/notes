package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"notes/internal/config"
	"notes/internal/lib/logger/sl"
	"notes/internal/notes/dto"

	// "notes/internal/notes/handlers"
	"notes/internal/notes/repository"
	"notes/internal/storage/postgresql"
	"notes/internal/utils"
	"os"
	"time"
)

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r)

			logger.Info(
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.statusCode),
				slog.Duration("time", time.Since(start)),
			)
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

	noteRepo := repository.NewNoteRepo(storage)
	ntd := &dto.NoteDTO{
		Title: utils.ToPtr("New note"),
	}
	_, err = noteRepo.Create(ntd)
	if err != nil {
		unw := errors.Unwrap(err)
		fmt.Println(unw)
		log.Printf("error: %v", err)
	}


	// mux := http.NewServeMux()

	// middlewareHandler := AuthorizationMiddleware("123")(mux)
	// middlewareHandler = loggingMiddleware(logger)(middlewareHandler)

	// server := http.Server{
	// 	Handler: middlewareHandler,
	// 	Addr:    "localhost:8083",
	// }

	// fmt.Printf("Server started on http://%s\n", server.Addr)

	// mux.HandleFunc("/", Hello)
	// mux.HandleFunc("/notes", handlers.CreateNote)

	// server.ListenAndServe()
}

func Hello(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello")
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
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

func AuthorizationMiddleware(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			formatedToken := fmt.Sprintf("Bearer %s", token)
			w.Header().Add("Authorization", formatedToken)
			next.ServeHTTP(w, r)
		})
	}
}
