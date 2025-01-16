package app

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"DeliveryTimePrediction/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-uuid"
	"github.com/jackc/pgx/v4"
)

type Storage interface {
	GetResult(ctx context.Context, id string) (float64, error)
}

type MessageQueue interface {
	Publish(topic string, data []byte) error
}

type App struct {
	storage Storage
	queue   MessageQueue
	topic   string
}

func New(storage Storage, queue MessageQueue, topic string) *App {
	return &App{
		storage: storage,
		queue:   queue,
		topic:   topic,
	}
}

func (a *App) GetResultHandler(w http.ResponseWriter, r *http.Request) {
	op := "GetResultHandler"
	ctx := r.Context()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	distance, err := a.storage.GetResult(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("%s. error getting result: %v", op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(distance)
	if err != nil {
		log.Printf("%s. error encoding result: %v", op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) PostTaskHandler(w http.ResponseWriter, r *http.Request) {
	op := "PostTaskHandler"
	var err error

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("%s. error parsing form: %v", op, err)
		a.sendError(w, err, http.StatusBadRequest)
		return
	}

	var task domain.Task
	task.Weather = r.Form.Get("weather")
	task.TrafficLevel = r.Form.Get("traffic_level")
	task.TimeOfDay = r.Form.Get("time_of_day")
	task.VehicleType = r.Form.Get("vehicle_type")
	task.Distance, err = strconv.ParseFloat(r.Form.Get("distance"), 64)
	if err != nil {
		log.Printf("%s. error parsing form, distance: %v", op, err)
		a.sendError(w, err, http.StatusBadRequest)
		return
	}
	task.PreparationTime, err = strconv.Atoi(r.Form.Get("preparation_time"))
	if err != nil {
		log.Printf("%s. error parsing form, preparation_time: %v", op, err)
		a.sendError(w, err, http.StatusBadRequest)
		return
	}
	task.CourierExperience, err = strconv.ParseFloat(r.Form.Get("courier_experience"), 64)
	if err != nil {
		log.Printf("%s. error parsing form, courier_experience: %v", op, err)
		a.sendError(w, err, http.StatusBadRequest)
		return
	}

	err = a.validateTask(task)
	if err != nil {
		log.Printf("%s. error validating task: %v", op, err)
		a.sendError(w, err, http.StatusBadRequest)
		return
	}

	task.ID, err = uuid.GenerateUUID()
	if err != nil {
		log.Printf("%s. error generating uuid: %v", op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("%s. task: %+v", op, task)

	data, err := json.Marshal(task)
	if err != nil {
		log.Printf("%s. error marshalling task: %v", op, err)
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}

	err = a.queue.Publish(a.topic, data)
	if err != nil {
		log.Printf("%s. error publishing task: %v", op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(task.ID)
	if err != nil {
		log.Printf("%s. error encoding task id: %v", op, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *App) sendError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(err)
	if err != nil {
		log.Printf("error encoding err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *App) validateTask(task domain.Task) error {
	switch task.Weather {
	case "Windy", "Clear", "Foggy", "Rainy", "Snowy":
	default:
		return errors.New("invalid weather")
	}
	switch task.TrafficLevel {
	case "Low", "Medium", "High":
	default:
		return errors.New("invalid traffic level")
	}

	switch task.TimeOfDay {
	case "Afternoon", "Evening", "Night", "Morning":
	default:
		return errors.New("invalid time of day")
	}

	switch task.VehicleType {
	case "Scooter", "Car", "Bike":
	default:
		return errors.New("invalid vehicle type")
	}

	switch {
	case task.PreparationTime < 0:
		return errors.New("invalid preparation time")
	case task.CourierExperience < 0:
		return errors.New("invalid courier experience")
	case task.Distance < 0:
		return errors.New("invalid distance")
	default:
	}

	return nil
}
