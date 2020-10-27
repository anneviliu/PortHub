package main

import (
	"errors"
	"fmt"
	"github.com/malfunkt/iprange"
	uuid "github.com/satori/go.uuid"
	"log"
	"net"
	"portScan_demo/forms"
	"strconv"
	"strings"
)

func CreateTaskID() string {
	id := uuid.NewV4()
	ids := id.String()
	return ids
}

// 解析ip 和 port的格式
func ResolveIPPortFormat(form *forms.PortScanForm) ([]net.IP, []int, error) {
	var ipList []net.IP
	list, err := iprange.ParseList(form.Ips)
	if err != nil {
		return nil, nil, errors.New("ip解析错误")
	}
	log.Printf("%+v", list)
	for _, ip := range list.Expand() {
		fmt.Println(ip)
		ipList = append(ipList, ip)
	}

	var ports []int
	portArr := strings.Split(strings.Trim(form.Ports, ","), ",")
	for _, v := range portArr {
		portArr2 := strings.Split(strings.Trim(v, "-"), "-")
		startPort, err := filterPort(portArr2[0])
		if err != nil {
			log.Println(err)
			return nil, nil, errors.New("port解析错误")
			//continue
		}
		ports = append(ports, startPort)
		if len(portArr2) > 1 {
			//添加第一个后面的所有端口
			endPort, _ := filterPort(portArr2[1])
			if endPort > startPort {
				for i := 1; i <= endPort-startPort; i++ {
					ports = append(ports, startPort+i)
				}
			}
		}
	}
	portList := Unique(ports)
	return ipList, portList, nil
}

func filterPort(str string) (int, error) {
	port, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	if port < 1 || port > 65535 {
		return 0, errors.New("Port out of range")
	}
	return port, nil
}

func Unique(arr []int) []int {
	var newArr []int
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
