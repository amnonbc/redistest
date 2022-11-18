package main

import (
	"context"
	"errors"
	"flag"
	"github.com/go-redis/redis/v8"
	"log"
	"sync/atomic"
	"time"
)

var (
	numMsgs = 0
	verbose = false
	nGot    int64
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.IntVar(&numMsgs, "n", 1000, "number of messages to send")
	flag.BoolVar(&verbose, "v", false, "vervose output")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go func() {
		nGot = 0
		pubsub := rdb.Subscribe(ctx, "ch")
		ch := pubsub.Channel()
		defer pubsub.Close()

		for range ch {
			n := atomic.AddInt64(&nGot, 1)
			if int(n) == numMsgs {
				cancel()
			}
		}

	}()
	start := time.Now()

	publish(rdb, ctx)

	<-ctx.Done()
	log.Println("got", atomic.LoadInt64(&nGot), "messages in", time.Since(start))

}

func publish(rdb *redis.Client, ctx context.Context) {

	for {
		err := rdb.Publish(ctx, "ch", "payload").Err()
		if errors.Is(err, context.Canceled) {
			break
		}
		if err != nil {
			panic(err)
		}

	}
}
