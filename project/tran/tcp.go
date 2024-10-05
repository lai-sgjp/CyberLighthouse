package tran

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// 一定要注意不要重复读取数据，并且UDP的server端一定要使用ReadFromUDP,WriteToUdp
func (t *Tcp) CreateSer(ports []string) {
	var (
		addrset []string
		input   string
	)
	scanner := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Which address do you want to listen to on the ?Please split them with a break.")
		fmt.Println("If you want to exit the pocess,please enter \"q\".")
		input, err := scanner.ReadString('\n')
		if err != nil {
			log.Println("Failed to read the addr.")
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			log.Println("You enter nothing.Please try again.")
			continue
		}
		if strings.ToLower(input) == "q" {
			log.Println("You have quitted the pocess")
			os.Exit(0)
		}
		break
	}

	addrset = strings.Split(input, " ")
	for _, port := range ports {
		//建立端口地址

		for _, addr := range addrset {
			addr, err := net.ResolveTCPAddr("tcp", addr+":"+port)
			if err != nil {
				log.Printf("Failed to get the correct TCP address %s: %v\n", port, err.Error())
				continue
			}

			listener, err := net.ListenTCP("tcp", addr)
			if err != nil {
				log.Printf("Failed to listen to the TCP address %s:%v\n", port, err.Error())
				continue
			}
			go func(listener *net.TCPListener) {
				for {
					conn, err := listener.Accept()
					if err != nil {
						log.Printf("Failed to accept the connection:%v\n", err.Error())
						continue
					}
					go t.Process(conn) //这个一定要放在for内部，因为变量生命周期的原因
				}

			}(listener)
		}

	}
}

func (t *Tcp) Process(conn interface{}) {
	realconn, _ := conn.(net.Conn)
	scanner := bufio.NewReader(realconn) //不同读取源就可以用不一样的
	for {
		//defer conn.Close()
		buf, err := scanner.ReadString('\n')
		if err != nil {
			log.Println("Failed to read the type and message sent by client:", err.Error())
			realconn.Close()
			return
		}
		buf = strings.TrimSpace(buf)
		typeinfo := string(buf[0]) //一个数字只占一个字节gangan
		switch typeinfo {
		case "0":
			err = t.textMode(realconn, buf[1:]) //所有字符串后面都会加上一个0结尾,而TrimSpace则是会将这个给去掉
			if err != nil {
				log.Printf("Failed to send message to the client:%v\n", err)
				realconn.Close()
				return
			}
		case "1":
			t.fileMode(realconn, scanner)
			realconn.Close()
			return
		default:
			log.Printf("Disapproval type:%s\n", typeinfo)
			//conn.Close()
		}
	}
}
