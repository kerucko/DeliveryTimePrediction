package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	app "DeliveryTimePrediction/internal/app"
	"DeliveryTimePrediction/internal/config"
	"DeliveryTimePrediction/internal/storage"
)

func main() {
	cfg := config.MustLoad()
	log.Printf("config: %+v", cfg)

	db, err := storage.New(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	service := app.New(db, nil)

	router.Post("/task", service.PostTaskHandler)
	router.Get("/check", service.GetResultHandler)

	log.Fatal(http.ListenAndServe("[::]:"+cfg.Server.Port, router))
}
