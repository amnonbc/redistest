package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// There is no error because go-redis automatically reconnects on error.
	pubsub := rdb.Subscribe(ctx, "mychannel1")

	// Close the subscription when we are done.
	defer pubsub.Close()
	ch := pubsub.Channel()

	for msg := range ch {
		log.Println(msg.Payload)
	}
}
