package database

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

var Redis redis.Conn

func init() {
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}

	Redis = conn
}