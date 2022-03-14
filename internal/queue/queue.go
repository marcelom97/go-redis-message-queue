package queue

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type Queue struct {
	client     *redis.Client
	streamName string
}

func NewQueue(client *redis.Client, streamName string) *Queue {
	return &Queue{client: client, streamName: streamName}
}

type QueueHandler interface {
	Publish(context.Context, string) error
	Consume(context.Context, string, string) ([]redis.XStream, error)
	Confirm(context.Context, string, string) error
	CreateConsumerGroup(context.Context, string) error
}

func (q Queue) Publish(ctx context.Context, message string) error {
	err := q.client.XAdd(ctx, &redis.XAddArgs{
		Stream: q.streamName,
		ID:     "*",
		Values: map[string]interface{}{
			"uuid":    uuid.NewString(),
			"message": message,
		},
	}).Err()
	return err
}

func (q Queue) Consume(ctx context.Context, consumerId string, consumersGroup string) ([]redis.XStream, error) {
	entries, err := q.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    consumersGroup,
		Consumer: consumerId,
		Streams:  []string{q.streamName, ">"},
		Count:    2,
		Block:    0,
		NoAck:    false,
	}).Result()
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (q Queue) Confirm(ctx context.Context, consumersGroup string, messageId string) error {
	err := q.client.XAck(ctx, q.streamName, consumersGroup, messageId).Err()
	return err
}

func (q Queue) CreateConsumerGroup(ctx context.Context, consumersGroup string) error {
	err := q.client.XGroupCreateMkStream(ctx, q.streamName, consumersGroup, "0").Err()
	return err
}
