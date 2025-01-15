package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"githb.com/Raunak9199/students-api/internal/config"
	"githb.com/Raunak9199/students-api/internal/handlers/students"
	"githb.com/Raunak9199/students-api/internal/storage/sqlite"
)

func main() {
	// load config

	cfg := config.MustLoad()

	// db setup

	storage, err := sqlite.New(cfg)

	if err != nil {
		log.Fatal(err)
	}
	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", students.New(storage))
	router.HandleFunc("GET /api/students/{id}", students.GetById(storage))

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started:", slog.String("address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to listen and serve. %s", err.Error())
		}
	}()

	<-done

	slog.Info("Shutting down the server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shut down the server.", slog.String("Error:", err.Error()))
	}

	slog.Info("Server shutdown completed successfully.")
}
