package tran

import (
	"log"
	"net"
)

//接口不允许在非本地类型上定义新方法，

func (u *Udp) CreateSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+port)
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
		go u.Process(conn)
	}

}

func (u *Udp) Process(conn interface{}) {
	realconn, _ := conn.(*net.UDPConn) //断言完成之后就可以用了

	//defer conn.Close()
	for {

		buf := make([]byte, 2048)
		n, clientAddr, err := realconn.ReadFromUDP(buf)
		log.Println(buf)
		if err != nil {
			log.Println("Failer to read the type and message sent by client:", err.Error())
			realconn.Close()
			return
		}
		typeinfo := string(buf[0])

		switch typeinfo {
		case "0":
			err = u.textMode(clientAddr, realconn, string(buf[1:n-1]))
			//log.Println("$%^")
			if err != nil {
				log.Printf("Failed to send message to the client:%v\n", err)
				realconn.Close()
				return
			}

		case "1":
			u.fileMode(realconn)
		default:
			log.Printf("Disapproval type:%s\n", typeinfo)
		}
	}
}
