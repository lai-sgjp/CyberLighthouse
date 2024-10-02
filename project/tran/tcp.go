package tran

import (
	"log"
	"net"
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
	defer conn.Close()
	buf := make([]byte, 8)

	n,err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to recieve type from the client:%v\n", err)
		conn.Close()
		return
	}

	typeinfo := string(buf[:n])
	switch typeinfo {
	case "0":
		newBuf := make([]byte,2048)
		for {
			err := textMode(conn,newBuf)
			if err != nil {
				break
			}
		}
	case "1":
		fileMode(conn)
	default:
		log.Printf("Disapproval type:%v\n",err)
		conn.Close() 
	}
	
}