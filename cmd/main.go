package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/marcelom97/go-redis-message-queue/producer"
	"github.com/marcelom97/go-redis-message-queue/queue"
)

type PingResponse struct {
	Message string `json:"message"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:        "redis:6379",
		DialTimeout: 20 * time.Second,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %s", err.Error())
	}

	streamName := "stream1"

	queue := queue.NewQueue(client, streamName)
	producer := producer.NewProducer(queue)
	producerHandler := NewProducerHandler(producer)

	r := mux.NewRouter()

	r.HandleFunc("/produce", producerHandler.Produce).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
