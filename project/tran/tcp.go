package tran
//被引用的包名必须和文件名一致
import (
	"fmt"
	"log"
	"net"
)

func CreateTCPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveTCPAddr("tcp", ":"+port)
		if err != nil {
			fmt.Printf("Failed to get the correct TCP address %s: %v", port, err)
			continue
		}
		//监听端口地址
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			fmt.Printf("Failed to listen to the TCP address %s:%v", port, err)
			continue
		}
		//没有错误表示正常监听（即监听已经建立）
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept the connection:%v", err)
			continue
		}
		go ProcessTCP(conn) //使用go采用多协程的方式从而并发处理
	}
}

func ProcessTCP(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)//这里的意思是将conn里面的读取到buf中
		if err != nil {
			log.Printf("Failed to recieve message from the client:%v", err)
			break
		}
		message := string(buf[:n])
		fmt.Printf("%s", message)

		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("Failed to send message to the client:%v", err)
			break
		}

	}
}