package main

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/marcelom97/go-redis-message-queue/internal/consumer"
	"github.com/marcelom97/go-redis-message-queue/internal/queue"
)

const (
	streamName     = "stream1"
	consumersGroup = "consumer-group-1"
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
		log.Fatal("Unable to connect to Redis: ", err)
	}

	log.Println("Connected to Redis server")

	q := queue.NewQueue(client, streamName)

	q.CreateConsumerGroup(ctx, consumersGroup)

	consumer := consumer.NewConsumer(q, consumersGroup)
	consumer.StartConsuming(ctx)
}
