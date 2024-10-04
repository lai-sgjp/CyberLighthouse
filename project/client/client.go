package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"CyberLighthouse/tran_c"
)

func main() {

	if len(os.Args) < 2 {
		tran_c.CreateConn() //未指定就默认模式
		return
	}

	protocolptr := flag.String("p", "tcp", "which protocol do you want to use(tcp/udp)")
	addrptr := flag.String("a", "127.0.0.1:8080", "which address do you want to send message") //a是go run的保留字
	modeptr := flag.String("m", "string", "which mode do you want to use(string/file)")        //感觉传整数和数字都一样。不如最后接收时转类型？
	contextptr := flag.String("c", "", "what do you want to send")

	flag.Parse()
	if *modeptr == "file" {
		for {
			if *contextptr == "" {
				log.Println("Please enter a filename.\nIf you want to quit the pocess,please enter \"Q\".")
				fmt.Scanf("%s", contextptr)
				if strings.ToUpper(*contextptr) == "Q" {
					log.Println("You have quitted the pocess.")
					return
				}
				continue
			}
			var conn net.Conn
			switch *protocolptr {
			case "tcp":
				conn = tran_c.CreateTCPConn(*addrptr)
			case "udp":
				conn = tran_c.CreateUDPConn(*addrptr)
			}
			tran_c.SendFile(conn, *contextptr)
			return
		}

	}
	if *protocolptr != "tcp" && *protocolptr != "udp" {
		log.Println("We don't support this service.We will use the \"tcp\" mode by deafult")
		tran_c.Choose("tcp", *addrptr, *contextptr)
		return
	}
	if *modeptr != "string" {
		log.Println("We don't support this type.We will use the \"string\" mode by deafult")
	}
	tran_c.Choose(*protocolptr, *addrptr, *contextptr)
	return
}
