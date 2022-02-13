package producer

import (
	"context"

	"github.com/marcelom97/go-redis-message-queue/queue"
)

type Producer struct {
	queue *queue.Queue
}

func NewProducer(queue *queue.Queue) *Producer {
	return &Producer{queue: queue}
}

func (p Producer) Produce(ctx context.Context, message string) error {
	err := p.queue.Publish(ctx, message)
	return err
}
