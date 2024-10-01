package tran_c
import (
	"bufio"
	"os"
	"strings"
	"fmt"
	"log"
	"net"
)

func sendMessage() (string, error) {
	//函数多个返回类型应该用括号包裹

	//var buf [2048]byte这个buff主要是用来接收信息而不是传入信息（因为知道大小）
	fmt.Println("What information do you want to send?Type Q to quit.")
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n')
	if err != nil {
		return "", err //直接返回一个空的字符串不就行了吗？!
	}
	trimmedInput := strings.Trim(input, "\n")
	return trimmedInput, nil
}

func CreateConn() {

	fmt.Println("Which port do you want to choose?")
	var port string
	fmt.Scanf("%s", &port)
	
	//暂且默认host为本机
	host := ""

	fmt.Println("Which way do you want to choose?Please enter 1 or 2.\n(1.udp\t2.tcp)")
	var choice int
	fmt.Scanf("%d", &choice)

	switch choice {
	case 2:
		/*
			var addr net.Addr
			var conn net.Conn
			var err error
		addr, err := net.ResolveTCPAddr("tcp", ":"+port)
		if err != nil {
			log.Fatal("Fail to get the correct TCP address:", err)
		}
		*/
		//conn, err := net.Dial("tcp", addr)//采用这种写法报错：cannot use addr (variable of type *net.TCPAddr) as string value in argument to net.Dial
		addr := fmt.Sprintf("%s:%s",host,port)
		conn ,err := net.Dial("tcp",addr)

		if err != nil {
			//fmt.Printf("Failed to dial:", err)
			log.Fatal("Failed to dial:", err)
		}
		
		defer conn.Close() //记得关闭连接

		for {

			trimmedInput, err := sendMessage()
			if err != nil {
				log.Fatal("Fail to input the infomation:", err)
			}
			if strings.ToUpper(trimmedInput) == "Q" {
				return
			}
			_, err = conn.Write([]byte(trimmedInput))
		}
	case 1:
		/*
		addr, err := net.ResolveUDPAddr("udp", ":"+port)
		if err != nil {
			log.Fatal("Fail to get the correct TCP address:", err)
		}
		*/
		addr := fmt.Sprintf("%s:%s",host,port)
		conn ,err := net.Dial("udp",addr)

		if err != nil {
			//fmt.Printf("Failed to dial:", err)//原来这是两个都会打印到屏幕上的
			log.Fatal("Failed to dial:", err)
		}
		defer conn.Close() //记得关闭连接

		for {

			trimmedInput, err := sendMessage() //一定要放在函数体内，否则就不能被修改变为无限循环
			if err != nil {
				log.Fatal("Fail to input the infomation:", err)
			}
			if strings.ToUpper(trimmedInput) == "Q" {
				return
			}
			_, err = conn.Write([]byte(trimmedInput))
		}
	default:
		//fmt.Println("Input Error!")
		log.Fatal("Input Error!")
	}
}