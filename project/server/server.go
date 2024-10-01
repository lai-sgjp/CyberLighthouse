package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func getport() ([]string, error) {
	fmt.Println("Please enter port seperate by space.(e.g.\"8000\" \"8001\" \"8002\")")
	//获得一连串的输入
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n')
	if err != nil {
		return nil, err
	}
	//格式化输入
	input = strings.TrimSpace(input)
	ports := strings.Split(input, " ")
	return ports, nil
}

func createUDPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveUDPAddr("udp", ":"+port)
		if err != nil {
			log.Printf("Failed to get the correct UDP address %s: %v", port, err)
			continue
		}

		//开始监听(或者说是连接)
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Printf("Failed to listen to UDP port %s:%v", port, err)
			continue
		}

		go processUDP(conn)
	}
}

func createTCPSer(ports []string) {
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
		go processTCP(conn) //使用go采用多协程的方式从而并发处理
	}
}

func processTCP(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
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

func processUDP(conn *net.UDPConn) { //是一个指针(UDP与TCP还是不一样的)
	//TCP则是强调传输的安全，怎么说就是用Accept确认一次并建立专门的通道/socket
	//UDP是只要传输包就可以了，所以只需要一个指向某个通道的指针
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Failed to receive message from the client:%v", err)
			break
		}

		fmt.Printf("%s", buf[:n])

		_, err = conn.WriteToUDP(buf[:n], clientAddr)
		if err != nil {
			log.Printf("Failed to send the message:%v", err)
		}
	}
}

func main() {
	ports, err := getport()
	if err != nil {
		log.Fatal("Failed to get the port:", err)
	}

	fmt.Println("Whether turn on the UDP service or not(y/n)")

	var choice rune //不是tune！！
	fmt.Scanln("%v", &choice)
	if choice == 'y' {
		createUDPSer(ports)
	}

	createTCPSer(ports)

	select {} //保持server端长时间不关闭
}
