package tran

import (
	"fmt"
	"log"
	"net"

)

func CreateUDPSer(ports []string) {
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

		go ProcessUDP(conn)
	}
}



func ProcessUDP(conn *net.UDPConn) { //是一个指针(UDP与TCP还是不一样的)
	//TCP则是强调传输的安全，怎么说就是用Accept确认一次并建立专门的通道/socket
	//UDP是只要传输包就可以了，所以只需要一个指向某个通道的指针、地址
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