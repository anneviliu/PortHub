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


	// Begin concurrent to scan
	ConLimit := make(chan int, form.Concurrent)
	// TODO:
	wg.Add(len(ports)*len(ips))
	for _, ip := range ips {
		_, err = conn.Do("SADD",ip,"running")
		if err != nil {
			log.Fatal(err)
		}

		if !scanner.ICMPRun(ip.String()) {
			_, err = conn.Do("SREM",ip,"running")
			if err != nil {
				log.Fatal(err)
			}
			continue
		}

		// Create a taskId for-each
		//_, err = conn.Do("SET",ip.String() + "taskid",CreateTaskID())
		//if err != nil {
		//	log.Fatal(err)
		//}

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

func GetSingleIpRes(c *gin.Context,ip string) {
	pool := database.NewPool()
	conn := pool.Get()
	port,err := redis.Strings(conn.Do("SMEMBERS",ip))
	if err != nil {
		log.Fatal(err)
	}
	for k,v := range port {
		if v == "running" {
			port = append(port[:k], port[k+1:]...)
			break
		}
	}

	c.JSON(200, forms.Response{StatusCode: 200, Messages: port, Data: map[string]interface{}{}})
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

	//for k,v := range ips {
	//	if strings.Contains(v, "taskid") {
	//		ips = append(ips[:k], ips[k+1:]...)
	//	}
	//}

	for _,ip := range ips {
		var taskId string

		port,err := redis.Strings(conn.Do("SMEMBERS",ip))
		if err!= nil {
			//log.Println(err)
			continue
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
		var long bool
		if len(port) > 10 {
			long = true
		} else {
			long = false
		}

		// Get taskID
		//isExist, err := redis.Bool(conn.Do("EXISTS",ip + "taskid"))
		//if isExist {
		//	taskId,err = redis.String(conn.Do("GET",ip + "taskid"))
		//	if err!= nil {
		//		log.Fatal(err)
		//	}
		//}

		tmp := map[string]interface{}{"port":port,"status":status[ip],"long":long,"taskId": taskId}
		data[ip] = tmp
	}

	c.JSON(200, forms.Response{StatusCode: 200, Messages: data, Data: map[string]interface{}{}})
}