package scanner

import (
	"fmt"
	"net"
	"sync"
)

var Timeout = 5
var Alive []string

func StartScanTask(ip net.IP, port int, wg *sync.WaitGroup) {
	//fmt.Println("Scan start")
	defer wg.Done()
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	//conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", tcpAddr.IP, tcpAddr.Port), time.Millisecond*time.Duration(Timeout))
	conn, err := net.DialTCP("tcp", nil, &tcpAddr)
	if err == nil {
		fmt.Println(ip, port, "is alive", conn, err)
		Alive = append(Alive, tcpAddr.IP.String()+string(tcpAddr.Port))
	}
	fmt.Println(err)
}
