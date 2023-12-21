package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/marcelom97/go-redis-message-queue/internal/queue"
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
	return c.queue.Confirm(ctx, c.consumersGroup, messageId)
}
