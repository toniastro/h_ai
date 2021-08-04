package handlers

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/gorilla/mux"
	"github.com/toniastro/h_ai/api/queues/rabbitMq"
	"log"
)

type Handler struct {
	RedisSync *redsync.Redsync
	Queue     *rabbitMq.Rabbit
}

func New(redis *redsync.Redsync, queue *rabbitMq.Rabbit) *Handler {
	return &Handler{
		RedisSync: redis,
		Queue:     queue,
	}
}

func (h *Handler) AddRoutes(route *mux.Router) {
	api := route.PathPrefix("/api").Subrouter()
	api.HandleFunc("/job", h.CreateJob).Methods("POST", "PUT")
	api.HandleFunc("/job/{id}", h.GetJobStatus).Methods("GET")
	log.Println("Loaded routes")
}
