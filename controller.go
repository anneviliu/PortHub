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
	"time"
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
	pool := database.NewPool()
	conn := pool.Get()

	c.JSON(200, forms.Response{StatusCode: 200, Messages: "", Data: map[string]interface{}{"taskId": CreateTaskID()}})

	ConLimit := make(chan int, form.Concurrent)

	wg.Add(len(ports)*len(ips) +1)
	for _, ip := range ips {
		_, err = conn.Do("SADD",ip,"running")
		if err != nil {
			log.Fatal(err)
		}
		for _, port := range ports {
			ConLimit <- 1
			go scanner.StartScanTask(ip, port, &wg,&ConLimit)
			RetResult()
		}
		time.Sleep(2000)
		_, err := conn.Do("SREM",ip,"running")
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(2000)
	}
	wg.Wait()
}

func RetResult() {
	pool := database.NewPool()
	conn := pool.Get()
	var infoArr []string
	for _, v := range scanner.Alive {
		infoArr = strings.Split(strings.Trim(v, ":"), ":")
		_, err := conn.Do("SADD",infoArr[0],infoArr[1])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetResult(c *gin.Context) {
	data := make(map[string]interface{})
	status := make(map[string]bool)

	pool := database.NewPool()
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

		//delete value "running" and return
		for k,v := range port {
			if v == "running" {
				status[ip] = true
				port = append(port[:k], port[k+1:]...)
				break
			}
			status[ip] = false
		}
		tmp := map[string]interface{}{"port":port,"status":status[ip]}
		data[ip] = tmp
	}

	c.JSON(200, forms.Response{StatusCode: 200, Messages: data, Data: map[string]interface{}{"taskId": CreateTaskID()}})
}