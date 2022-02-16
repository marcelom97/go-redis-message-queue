package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
	w.Header().Add("Content-Type", "application/json")
	ctx := r.Context()

	var b PublishRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&b)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, fmt.Sprintf("Bad Request. Wrong Type provided for field: %s", unmarshalErr.Field), http.StatusBadRequest)
			return
		}
		errorResponse(w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest)
		return
	}

	err = h.producer.Produce(ctx, b.Message)
	if err != nil {
		errorResponse(w, fmt.Sprintf("Something went wrong: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(PublishRequest{
		Message: "Message was published successfully!",
	})
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	json.NewEncoder(w).Encode(PublishErrorResponse{Error: message})
}
