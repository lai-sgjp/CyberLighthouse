package tran_c
import (
	"fmt"
	"log"
	"net"
)

type quitError struct {
	message string
}
func (e quitError) Error() string{
	return e.message
}

//创立指定连接TCP
func CreateTCPConn(addr string) net.Conn {
	conn ,err := net.Dial("tcp",addr)

	if err != nil {
		//fmt.Printf("Failed to dial:", err)
		log.Fatal("Failed to dial:", err)
	}
		
	
	return conn
}
//创立指定连接UDP
func CreateUDPConn(addr string) net.Conn {
	conn ,err := net.Dial("udp",addr)

	if err != nil {
		log.Fatal("Failed to dial:", err)
	}
	return conn
}
//选择某一个指定的连接(默认模式下)
//文件模式在client.go中已经定义
func Choose(tran string,addr string,message string) {

	switch tran {
	case "tcp":
		conn := CreateTCPConn(addr)
		defer conn.Close()
		//message,err := sendMessage()
		//if err != nil {
		//	log.Fatal("Failed to read the message:",err)
		//}
		err := defaultSend(conn)
		if err != nil {
			break
		}
	case "udp":
		conn := CreateUDPConn(addr)
		defer conn.Close()
		err := defaultSend(conn)
		if err != nil {
			break
		}
	}
}


//发送信息(默认)因为是循环的，所以难以拆分（或者说没思路）
func CreateConn() {

	fmt.Println("Which port do you want to choose?")
	var port string
	fmt.Scanf("%s", &port)
	
	//暂且默认host为本机
	host := ""

	fmt.Println("Which way do you want to choose?Please enter 1 or 2.\n(1.udp\t2.tcp)")
	var choice int
	fmt.Scanf("%d", &choice)

	switch choice {
	case 2:
		addr := fmt.Sprintf("%s:%s",host,port)
		
		conn := CreateTCPConn(addr)
		defer conn.Close() //记得关闭连接
		for {
			err := defaultSend(conn)
			if err != nil {
				break
			}
		}
		
	case 1:
		addr := fmt.Sprintf("%s:%s",host,port)
		conn := CreateUDPConn(addr)
		
		defer conn.Close() 

		for {
			err := defaultSend(conn)
			if err != nil {
				break
			}
		}
		
	default:
		log.Fatal("Input Error!")
	}
}