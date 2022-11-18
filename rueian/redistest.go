package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/rueian/rueidis"
	"log"
	"sync/atomic"
	"time"
)

var (
	numMsgs = 0
	verbose = false
	delay   = time.Microsecond
)

func main() {
	flag.IntVar(&numMsgs, "n", 1000, "number of messages to send")
	flag.BoolVar(&verbose, "v", false, "vervose output")
	flag.DurationVar(&delay, "delay", time.Microsecond, "how long to wait")
	flag.Parse()
	c, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	fmt.Println(got, err)
}

func publisher(ctx context.Context, c rueidis.Client) {
	tick := time.NewTicker(delay)
	defer tick.Stop()

	for range tick.C {
		err := c.Do(ctx, c.B().Publish().Channel("ch").Message("msg").Build()).Error()
		if errors.Is(err, context.Canceled) {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}
