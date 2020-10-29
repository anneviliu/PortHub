package main

import (
	"PortHub/database"
	"PortHub/forms"
	"PortHub/scanner"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"strings"
	"sync"
)


func ScannerController(c *gin.Context) {
	form := new(forms.PortScanForm)
	err := c.BindJSON(form)
	if err != nil {
		c.JSON(400, forms.Response{StatusCode: 400, Messages: err.Error(), Data: nil})
		return
	}

	ips, ports, err := ResolveIPPortFormat(form)
	if err != nil {
		c.JSON(400, forms.Response{StatusCode: 400, Messages: err.Error(), Data: nil})
		return
	}
	var wg sync.WaitGroup
	c.JSON(200, forms.Response{StatusCode: 200, Messages: "", Data: map[string]interface{}{"taskId": CreateTaskID()}})

	ConLimit := make(chan int, form.Concurrent)

	wg.Add(len(ports)*len(ips) +1)
	for _, ip := range ips {
		for _, port := range ports {
			ConLimit <- 1
			go scanner.StartScanTask(ip, port, &wg,&ConLimit)
			RetResult()
		}
		_, err := database.Redis.Do("SREM",ip,"running")
		if err != nil {
			log.Fatal(err)
		}
	}

	wg.Wait()
}

func RetResult() {
	var infoArr []string
	for _, v := range scanner.Alive {
		infoArr = strings.Split(strings.Trim(v, ":"), ":")
		_, err := database.Redis.Do("SADD",infoArr[0],infoArr[1])
		if err != nil {
			log.Fatal(err)
		}

		_, err = database.Redis.Do("SADD",infoArr[0],"running")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetResult(c *gin.Context) {
	data := make(map[string][]string)
	pool := NewPool()
	conn := pool.Get()
	defer conn.Close()

	ips, err := redis.Strings(conn.Do("keys","*"))
	if err != nil {
		log.Fatal(err)
	}

	for _,ip := range ips {
		port,err := redis.Strings(conn.Do("SMEMBERS",ip))
		if err!= nil {
			log.Fatal(err)
		}

		// delete value "running"
		for k,v := range port {
			if v == "running" {
				port = append(port[:k], port[k+1:]...)
				break
			}
		}
		data[ip] = port
	}

	c.JSON(200, forms.Response{StatusCode: 200, Messages: data, Data: map[string]interface{}{"taskId": CreateTaskID()}})
}