package forms

import "github.com/jinzhu/gorm"

type PortScanForm struct {
	Ips   string `json:"ips"`
	Ports string `json:"ports"`
}

type Response struct {
	StatusCode int         `json:"statusCode" example:"200"`
	Messages   string      `json:"messages" example:"错误信息"`
	Data       interface{} `json:"data"`
}

type ScannerDb struct {
	gorm.Model
	Ip     string `gorm:"ip"`
	Ports  string `gorm:"ports"`
	TaskId string `gorm:"taskId"`
}
