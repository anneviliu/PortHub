package main

import (
	"PortHub/forms"
	"PortHub/scanner"
	"github.com/gin-gonic/gin"
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

	// 暂时先实现单ip多端口的扫描
	// 从ip port表达式中解析出单个ip+port
	ips, ports, err := ResolveIPPortFormat(form)
	if err != nil {
		c.JSON(400, forms.Response{StatusCode: 400, Messages: err.Error(), Data: nil})
		return
	}
	var wg sync.WaitGroup

	//wg.Add(len(ports))
	for _, ip := range ips {
		for _, port := range ports {
			wg.Add(1)
			go scanner.StartScanTask(ip, port, &wg)
			time.Sleep(1000)
		}
	}
	wg.Wait()
	//fmt.Println(scanner.Alive)
	c.JSON(200, forms.Response{StatusCode: 200, Messages: "", Data: map[string]interface{}{"taskId": CreateTaskID()}})
}
