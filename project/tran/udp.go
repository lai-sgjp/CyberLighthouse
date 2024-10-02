package tran

import (
	"log"
	"net"
)

func CreateUDPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveUDPAddr("udp", ":"+port)
		if err != nil {
			log.Printf("Failed to get the correct UDP address %s: %v\n", port, err)
			continue
		}

		//开始监听(或者说是连接)
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Printf("Failed to listen to UDP port %s:%v\n", port, err)
			continue
		}

		go ProcessUDP(conn)
	}
}

func ProcessUDP(conn *net.UDPConn) { 

	defer conn.Close()
	buf := make([]byte, 8)
	n,_,err := conn.ReadFromUDP(buf)
	typeinfo := string(buf[:n])
	if err != nil {
		log.Printf("Failed to recieve type from the client:%v\n", err)
		conn.Close() 
		return
	}
	switch typeinfo {
	case "0":
		newBuf := make([]byte,2048)
		for {
			err := textModeudp(conn,newBuf)
			if err != nil {
				break
			}
		}
	case "1":
		for {
			fileModeudp(conn)
		}
	default:
		log.Printf("Disapproval type:%v\n",err)
		conn.Close() 
	}
}