package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	app "DeliveryTimePrediction/internal/app"
	"DeliveryTimePrediction/internal/config"
	"DeliveryTimePrediction/internal/kafka/consumer"
	"DeliveryTimePrediction/internal/kafka/producer"
	"DeliveryTimePrediction/internal/storage"
)

func main() {
	cfg := config.MustLoad()
	log.Printf("config: %+v", cfg)

	ctx := context.Background()
	wg := &sync.WaitGroup{}

	db, err := storage.New(context.Background(), cfg.Postgres)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	kafka_producer, err := producer.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatalf("failed to initialize kafka producer: %v", err)
	}

	consumer_handler := consumer.NewConsumerGroupHandler(db)
	kafka_consumer, err := consumer.NewConsumerGroup(
		cfg.Kafka.Brokers,
		cfg.Kafka.ConsumerGroup,
		[]string{cfg.Kafka.Topics.Completed},
		consumer_handler,
	)
	if err != nil {
		log.Fatalf("failed to initialize kafka consumer: %v", err)
	}
	defer kafka_consumer.Close()

	kafka_consumer.Run(ctx, wg)

	service := app.New(db, kafka_producer, cfg.Kafka.Topics.Tasks)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/task", service.PostTaskHandler)
	router.Get("/check/{id}", service.GetResultHandler)

	log.Fatal(http.ListenAndServe("[::]:"+cfg.Server.Port, router))
	wg.Wait()
}
