package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"CyberLighthouse/dns/process"
	"CyberLighthouse/tran_c"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Which addr do you want to choose?")
		var addr string
		fmt.Scanf("%s", &addr)
		if addr == "" {
			log.Println("You enter none to the buffer.We will use \"127.0.0.1\" as default")
			addr = "127.0.0.1"
		}

		//暂且默认host为本机
		fmt.Println("Which port do you want to choose?")
		var port string
		fmt.Scanf("%s", &port)
		if port == "" {
			log.Println("You enter none to the buffer.We will use \"8080\" as default")
			port = "8080"
		}

		fmt.Println("Which way do you want to choose?Please enter \"udp\" or \"tcp\"")
		var choice string
		fmt.Scanf("%s", &choice)
		addr = fmt.Sprintf("%s:%s", addr, port)
		switch choice {
		case "tcp":
			tran_c.Choose("tcp", addr, "")
		case "udp":
			tran_c.Choose("udp", addr, "")
		default:
			log.Println("You enter protocol we don't support..We will use \"tcp\" as default")
			tran_c.Choose("tcp", addr, "")
		} //未指定就默认模式
		return
	}

	protocolptr := flag.String("p", "tcp", "which protocol do you want to use(tcp/udp)")
	addrptr := flag.String("a", "127.0.0.1:8080", "which address do you want to send message/ask for the DNS server") //a是go run的保留字
	modeptr := flag.String("m", "string", "which mode do you want to use(string/file)")                               //感觉传整数和数字都一样。不如最后接收时转类型？
	contextptr := flag.String("c", "", "what do you want to send/which domain you want to ask")
	dnsptr := flag.Bool("dns", false, "whether ask for the DNS server or not(true/false).If you choose \"true\",you can define \"-a\" ,\"-p\" and \"-c\"")

	flag.Parse()
	if *dnsptr {
		if (strings.Replace(*addrptr, " ", "", -1) == "") || (*addrptr == "127.0.0.1:8080") {
			fmt.Println("Please enter a DNS server and the port(port usually is 53):") //期待加上超时
			fmt.Scanf("%s", addrptr)
			if (strings.Replace(*addrptr, " ", "", -1) == "") || (*addrptr == "127.0.0.1:8080") {
				log.Println("Since you enter nothing/break, we will use 8.8.8.8:53 by default.")
				*addrptr = "8.8.8.8:53"
			}
		}
		if *contextptr == "" {
			fmt.Println("Please enter which domain address you want to analyse:")
			fmt.Scanf("%s", contextptr) //这个地方一个是输入Println后面会带一个\n，而且输入的存储块应该在*contextptr里面

			if strings.Replace(*addrptr, " ", "", -1) == "" {
				log.Println("Since you enter nothing/break, we will exit the pocess.")
				os.Exit(1)
			}
		}

		if *protocolptr != "tcp" && *protocolptr != "udp" {
			fmt.Println("Which way do you want to choose?Please enter \"udp\" or \"tcp\"")
			fmt.Scanf("%s", protocolptr)
			if *protocolptr != "tcp" && *protocolptr != "udp" {
				log.Println("You enter protocol we don't support..We will use \"udp\" as default")
				*protocolptr = "udp"
			} //未指定就默认模式
		}
		//u := tran_c.Udp{}
		for i := 0; i < 4; i++ {
			e := process.DNS(*addrptr, *contextptr, *protocolptr)
			if e {
				return
			}
		}
		fmt.Println("Failed to get the response for a long time...")
		return
	}
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
			switch *protocolptr {
			case "tcp":
				t := tran_c.Tcp{}
				conn := t.CreateConn(*addrptr)
				t.SendFile(conn, *contextptr)
			case "udp":
				u := tran_c.Udp{}
				_, conn := u.CreateConn(*addrptr)
				u.SendFile(conn, *contextptr) //想着能不能用接口进行类型断言然后使用？
			}

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
}

//整体抽象出来，启动客户端
