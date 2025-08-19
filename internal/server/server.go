package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"iron-insight/internal/database"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	slog.Debug("Debug message")

	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug, // Show debug logs
	// }))

	logger := slog.New(slog.NewJSONHandler(&PrettyJSONWriter{w: os.Stdout}, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	NewServer := &Server{
		port: port,

		db: database.NewService(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
