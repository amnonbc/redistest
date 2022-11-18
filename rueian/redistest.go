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
	numMsgs    = 0
	verbose    = false
	batchSize  = 1
	concurrent = false
)

func main() {
	flag.IntVar(&numMsgs, "n", 10000, "number of messages to send")
	flag.IntVar(&batchSize, "batchsize", 100, "size of batch to publish")
	flag.BoolVar(&concurrent, "concurrent", false, "do each publish in a goroutine")
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
	if concurrent {
		go concurrentPublisher(ctx, c)
	} else if batchSize > 1 {
		go batchPublisher(ctx, c)
	} else {
		go publisher(ctx, c)
	}
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

func batchPublisher(ctx context.Context, c rueidis.Client) {
	msgs := make(rueidis.Commands, batchSize)
	for n := 0; n < numMsgs; n += batchSize {
		for i := range msgs {
			msgs[i] = c.B().Publish().Channel("ch").Message("msg").Build()
		}
		for _, resp := range c.DoMulti(ctx, msgs...) {
			err := resp.Error()
			if err != nil {
				log.Println("error publishing", err)
			}
		}
	}
}

func concurrentPublisher(ctx context.Context, c rueidis.Client) {
	for i := 0; i < numMsgs; i++ {
		go func() {
			err := c.Do(ctx, c.B().Publish().Channel("ch").Message("msg").Build()).Error()
			if errors.Is(err, context.Canceled) {
				return
			}
			if err != nil {
				log.Println("publish", err)
			}
		}()
	}
}
