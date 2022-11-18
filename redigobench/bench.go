package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalln(err)

	}
	defer c.Close()
	c.Send("SUBSCRIBE", "mychannel1")
	c.Flush()
	for {
		reply, err := c.Receive()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(reply)
		// process pushed message
	}
}
