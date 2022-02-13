package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/marcelom97/go-redis-message-queue/queue"
)

type Consumer struct {
	queue          *queue.Queue
	id             string
	consumersGroup string
}

func NewConsumer(queue *queue.Queue, consumersGroup string) *Consumer {
	return &Consumer{
		queue:          queue,
		id:             uuid.NewString(),
		consumersGroup: consumersGroup,
	}
}

type ConsumerHandler interface {
	StartConsuming(context.Context) error
	Consumed(string) error
}

func (c Consumer) StartConsuming(ctx context.Context) error {
	for {
		entries, err := c.queue.Consume(ctx, c.id, c.consumersGroup)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(entries[0].Messages); i++ {
			messageId := entries[0].Messages[i].ID
			values := entries[0].Messages[i].Values
			emailUuid := fmt.Sprintf("%v", values["uuid"])
			emailMessage := fmt.Sprintf("%v", values["message"])
			fmt.Printf("consumer: [%s] uuid: [%s] message: [%s]\n", c.id, emailUuid, emailMessage)
			c.Consumed(ctx, messageId)
		}
	}
}

func (c Consumer) Consumed(ctx context.Context, messageId string) error {
	err := c.queue.Confirm(ctx, c.consumersGroup, messageId)
	return err
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
		log.Fatal("Unable to connect to Redis: ", err)
	}

	log.Println("Connected to Redis server")

	streamName := "stream1"
	consumersGroup := "consumer-group-1"
	q := queue.NewQueue(client, streamName)

	q.CreateConsumerGroup(ctx, consumersGroup)

	consumer := NewConsumer(q, consumersGroup)
	consumer.StartConsuming(ctx)
}
