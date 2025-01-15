package cmd

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	app "DeliveryTimePrediction/internal/app"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	service := app.New(nil)

	router.Post("/task", service.PostTaskHandler)
	router.Get("/check", service.GetResultHandler)
}
