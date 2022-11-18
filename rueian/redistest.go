package main

import (
	"context"
	"errors"
	"flag"
	"github.com/rueian/rueidis"
	"log"
	"sync/atomic"
	"time"
)

var (
	numMsgs = 0
	verbose = false
)

func main() {
	flag.IntVar(&numMsgs, "n", 10000, "number of messages to send")
	flag.BoolVar(&verbose, "v", false, "vervose output")
	flag.Parse()
	c, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := time.Now()
	go publisher(ctx, c)
	var got int64
	err = c.Receive(ctx, c.B().Subscribe().Channel("ch").Build(), func(msg rueidis.PubSubMessage) {
		n := atomic.AddInt64(&got, 1)
		if verbose {
			log.Println("got msg", n, msg.Channel, msg.Message)
		}
		if int(n) == numMsgs {
			cancel()
		}
	})
	if !errors.Is(err, context.Canceled) {
		log.Println(err)
	}
	log.Println("got", atomic.LoadInt64(&got), "messages in", time.Since(start))
}

func publisher(ctx context.Context, c rueidis.Client) {
	for {
		err := c.Do(ctx, c.B().Publish().Channel("ch").Message("msg").Build()).Error()
		if errors.Is(err, context.Canceled) {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}
