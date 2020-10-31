package main

import (
	"PortHub/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	database.InitDb()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	
	r.POST("/api/createPortScanTask", func(c *gin.Context) {
		ScannerController(c)
	})
	
	r.GET("/result", func(c *gin.Context) {
		c.HTML(http.StatusOK, "result.html", gin.H{})
	})
	
	r.GET("/api/getResult", func(c *gin.Context) {
		GetResult(c)
	})
	r.Run()

	//var wg sync.WaitGroup
	//s := make(chan os.Signal)
	////wg.Add(1)
	//go func() {
	//	signal.Notify(s)
	//	t:= <-s
	//	fmt.Println(t)
	//	//wg.Done()
	//}()
	//wg.Wait()
}