package main

import (
	"encoding/json"
	"net/http"

	"github.com/marcelom97/go-redis-message-queue/producer"
)

type ProducerHandler struct {
	producer *producer.Producer
}

func NewProducerHandler(producer *producer.Producer) *ProducerHandler {
	return &ProducerHandler{producer: producer}
}

type PublishRequest struct {
	Message string `json:"message"`
}

type PublishErrorResponse struct {
	Error string `json:"error"`
}

type PublishResponse struct {
	Message string `json:"message"`
}

func (h ProducerHandler) Produce(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var b PublishRequest

	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.producer.Produce(ctx, b.Message)
	if err != nil {
		json.NewEncoder(w).Encode(PublishErrorResponse{
			Error: "something went wrong",
		})
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PublishRequest{
		Message: "Message was published successfully!",
	})
}
