package scanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

var Alive []string

func StartScanTask(ip net.IP, port int, wg *sync.WaitGroup,ConLimit *chan int) {
	defer wg.Done()
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	conn, _ := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", tcpAddr.IP, tcpAddr.Port), time.Millisecond*time.Duration(500))
	if conn != nil {
		fmt.Println(ip, port, "is alive")
		err := conn.Close()
		if err != nil {
			panic(err)
		}

		Alive = append(Alive, tcpAddr.IP.String()+":" +strconv.Itoa(tcpAddr.Port))
	}

	<- *ConLimit
}
