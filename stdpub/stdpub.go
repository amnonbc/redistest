package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func main() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	defer rdb.Close()
	rdb.Publish(ctx, "mychannel1", "msg from stdpub: "+time.Now().String())
}
