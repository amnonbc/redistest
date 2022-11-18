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
	err = c.Send("SUBSCRIBE", "ch")
	if err != nil {
		log.Fatalln(err)
	}
	err = c.Flush()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		reply, err := c.Receive()
		if err != nil {
			log.Fatalln(err)
		}
		arr, ok := reply.([]interface{})
		if !ok {
			continue
		}

		for _, w := range arr {
			_, ok := w.([]byte)
			pat := "%v "
			if ok {
				pat = "%s "
			}
			fmt.Printf(pat, w)
		}
		fmt.Println()
	}
}
