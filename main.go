package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"PortHub/database"
)

func main() {
	database.InitDb()
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	
	r.POST("/createPortScanTask", func(c *gin.Context) {
		ScannerController(c)
	})
	
	r.GET("/result", func(c *gin.Context) {
		c.HTML(http.StatusOK, "result.html", gin.H{})
	})
	
	r.GET("/api/getResult", func(c *gin.Context) {
		GetResult(c)
	})
	r.Run()
}