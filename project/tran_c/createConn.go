package tran_c

import (
	"fmt"
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
func CreateTCPConn(addr string) net.Conn {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatal("Failed to dial:", err)
	}

	return conn
}

// 创立指定连接UDP
func CreateUDPConn(addr string) net.Conn {
	conn, err := net.Dial("udp", addr)

	if err != nil {
		log.Fatal("Failed to dial:", err)
	}
	return conn
}

// 选择某一个指定的连接(默认模式下)
// 文件模式在client.go，send中已经定义
func Choose(tran string, addr string, message string) {
	switch tran {
	case "tcp":
		conn := CreateTCPConn(addr)
		defer conn.Close()

		for i := 0; i < 10; i++ {
			typeinfo := fmt.Sprintf("%d\n", 0)
			_, err := conn.Write([]byte(typeinfo))
			if err != nil {
				log.Println("Failed to send the type", err.Error())
				return
			}
			if message == "" {
				message, err = Message()
				if err != nil {
					log.Println("Fail to get the message...")
					continue
				}
			}
			err = defaultSend(conn, message)
			if err != nil {
				log.Println("Failed to send the message.", err.Error())
				break
			}
			message = ""
		}
	case "udp":
		conn := CreateUDPConn(addr)
		defer conn.Close()
		for i := 0; i < 10; i++ {
			typeinfo := fmt.Sprintf("%d", 0)
			_, err := conn.Write([]byte(typeinfo))
			if err != nil {
				log.Println("Failed to send the type", err.Error())
				return
			}
			if message == "" {
				message, err = Message()
				if err != nil {
					log.Println("Fail to get the message...")
					continue
				}
			}
			err = defaultSend(conn, message)
			if err != nil {
				log.Println("Faile to send the message.", err.Error())
				break
			}
		}
	}
}

// 发送信息(默认)因为是循环的，所以难以拆分（或者说没思路）
func CreateConn() {
	fmt.Println("Which addr do you want to choose?")
	var addr string
	fmt.Scanf("%s", &addr)
	//暂且默认host为本机
	fmt.Println("Which port do you want to choose?")
	var port string
	fmt.Scanf("%s", &port)
	if addr == "" || port == "" {
		log.Println("You enter none to the buffer.We will use \"127.0.0.1:8080\" as default")
	}

	fmt.Println("Which way do you want to choose?Please enter \"udp\" or \"tcp\"")
	var choice string
	fmt.Scanf("%s", &choice)
	addr = fmt.Sprintf("%s:%s", addr, port)
	switch choice {
	case "tcp":
		Choose("tcp", addr, "")
	case "udp":
		Choose("udp", addr, "")
	default:
		log.Println("You enter protocol we don't support..We will use \"tcp\" as default")
		Choose("tcp", addr, "")
	}
}
