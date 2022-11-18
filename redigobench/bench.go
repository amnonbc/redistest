package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalln(err)

	}
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe("ch")
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			log.Fatalln(v)
		}
	}

}
