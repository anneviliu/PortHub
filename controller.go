package main

import (
	"PortHub/database"
	"PortHub/forms"
	"PortHub/scanner"
	"github.com/gin-gonic/gin"
	"log"
	"sync"
)

func ScannerController(c *gin.Context) {
	form := new(forms.PortScanForm)
	err := c.BindJSON(form)
	if err != nil {
		c.JSON(400, forms.Response{StatusCode: 400, Messages: err.Error(), Data: nil})
		return
	}

	// 暂时先实现单ip多端口的扫描
	// 从ip port表达式中解析出单个ip+port
	ips, ports, err := ResolveIPPortFormat(form)
	if err != nil {
		c.JSON(400, forms.Response{StatusCode: 400, Messages: err.Error(), Data: nil})
		return
	}
	var wg sync.WaitGroup
	c.JSON(200, forms.Response{StatusCode: 200, Messages: "", Data: map[string]interface{}{"taskId": CreateTaskID()}})

	// Cool concurrent count
	ConLimit := make(chan int, form.Concurrent)

	wg.Add(len(ports)*len(ips) +1 )
	for _, ip := range ips {
		for _, port := range ports {
			ConLimit <- 1
			go scanner.StartScanTask(ip, port, &wg,&ConLimit)
		}
	}
	wg.Wait()
}

func RetResult(c *gin.Context) {

	_, err := database.Redis.Do("SADD",ip,port)
	if err != nil {
		log.Fatal(err)
	}
}
