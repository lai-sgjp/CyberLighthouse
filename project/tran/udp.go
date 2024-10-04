package tran

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func CreateUDPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveUDPAddr("udp", ":"+port)
		if err != nil {
			log.Printf("Failed to get the correct UDP address %s: %v\n", port, err.Error())
			continue
		}

		//开始监听(或者说是连接)
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Printf("Failed to listen to UDP port %s:%v\n", port, err.Error())
			continue
		}
		go ProcessUDP(conn)
	}

}

func ProcessUDP(conn *net.UDPConn) {
	scanner := bufio.NewReader(conn)
	//defer conn.Close()
	for {
		buf, err := scanner.ReadString('\n')
		if err != nil {
			log.Println("Failer to read the type and message sent by client:", err.Error())
			conn.Close()
			return
		}
		typeinfo := string(buf[0])
		switch typeinfo {
		case "0":
			for {
				buf, err := scanner.ReadString('\n')
				buf = strings.TrimSpace(buf)
				err = textModeudp(conn, buf[:len(buf)-1])
				if err != nil {
					log.Println("Failed to read the message:", err.Error())
					conn.Close()
					return
				}
			}
		case "1":
			fileModeudp(conn)
		default:
			log.Printf("Disapproval type:%c\n", typeinfo)
		}
	}
}
