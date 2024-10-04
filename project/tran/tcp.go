package tran

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func CreateTCPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)
		if err != nil {
			log.Printf("Failed to get the correct TCP address %s: %v\n", port, err)
			continue
		}

		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Printf("Failed to listen to the TCP address %s:%v\n", port, err)
			continue
		}
		go func(listener *net.TCPListener) {
			for {
				conn, err := listener.Accept()
				if err != nil {
					log.Printf("Failed to accept the connection:%v\n", err)
					continue
				}
				go ProcessTCP(conn) //这个一定要放在for内部，因为变量生命周期的原因
			}

		}(listener)

	}
}

func ProcessTCP(conn net.Conn) {
	scanner := bufio.NewReader(conn)
	for {
		//defer conn.Close()
		buf, err := scanner.ReadString('\n')
		if err != nil {
			log.Println("Failed to read the type and message sent by client:", err.Error())
			return
		}
		typeinfo := string(buf[0]) //一个数字只占一个字节gangan
		switch typeinfo {
		case "0":
			for {
				//newBuf := make([]byte, 2048)
				buf, err := scanner.ReadString('\n')
				buf = strings.TrimSpace(buf)
				err = textMode(conn, buf[:len(buf)-1]) //所有字符串后面都会加上一个0结尾
				if err != nil {
					log.Println("Failed to read message:", err)
					return
				}
			}

		case "1":
			fileMode(conn)
		default:
			log.Printf("Disapproval type:%c\n", typeinfo)
			//conn.Close()
		}
	}
}
