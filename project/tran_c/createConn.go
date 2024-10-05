package tran_c

import (
	"log"
	"net"
)

// 自定义错误
type quitError struct {
	message string
}

func (e quitError) Error() string {
	return e.message
}

// 创立指定连接TCP
func (t *Tcp) CreateConn(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal("Failed to dial:", err)
	}

	return conn
}

// 创立指定连接UDP
func (u *Udp) CreateConn(addr string) (*net.UDPAddr, *net.UDPConn) {
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Println("Failed to resolve the UDP address:", err.Error())
	}
	conn, err := net.DialUDP("udp", nil, serverAddr) //UDP一定要使用专用

	if err != nil {
		log.Fatal("Failed to dial:", err)
	}
	return serverAddr, conn
}

// 选择某一个指定的连接(默认模式下)
// 文件模式在client.go，send中已经定义
func Choose(tran string, addr string, message string) {
	switch tran {
	case "tcp":
		t := Tcp{}
		conn := t.CreateConn(addr)
		defer conn.Close()
		var err error
		for i := 0; i < 10; i++ {
			message, err = getmessage(message)
			if err != nil {
				log.Println("Fail to get the message...")
				continue
			}

			err = t.defaultSend(conn, message)
			if err != nil {
				log.Println("Failed to send the message:", err.Error())
				continue
			}
			message = ""
		}

	case "udp":
		u := Udp{}
		_, conn := u.CreateConn(addr)
		defer conn.Close()
		var err error
		for i := 0; i < 10; i++ {
			message, err = getmessage(message)
			if err != nil {
				log.Println("Fail to get the message...")
				continue
			}

			err = u.defaultSend(conn, message)
			if err != nil {
				log.Println("Failed to send the message.", err.Error())
				break
			}
			message = ""
		}
	}
}
