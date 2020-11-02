package database

import (
	"github.com/gomodule/redigo/redis"
)

var Redis redis.Conn

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