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

func NewPoolFunc()(redis.Conn, error){
	return redis.Dial("tcp", ":6379")
}

func NewPool()(* redis.Pool){
	return &redis.Pool{
		MaxIdle: 50,
		Dial: NewPoolFunc,
		Wait: true,
	}
}