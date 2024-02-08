package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/marcelom97/go-redis-message-queue/internal/handlers"
	"github.com/marcelom97/go-redis-message-queue/internal/producer"
	"github.com/marcelom97/go-redis-message-queue/internal/queue"
)

const (
	streamName = "stream1"
)

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

	queue := queue.NewQueue(client, streamName)
	producer := producer.NewProducer(queue)
	producerHandler := handlers.NewProducerHandler(producer)

	r := http.NewServeMux()

	r.HandleFunc("POST /produce", producerHandler.Produce)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
