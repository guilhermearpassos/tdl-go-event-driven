package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewStdLogger(false, false)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	messages, err := sub.Subscribe(context.Background(), "progress")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		level := string(msg.Payload)
		fmt.Printf("Message ID: %v - %v%%", msg.UUID, level)
		msg.Ack()
	}
}
