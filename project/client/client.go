package main 

import (
	"http/net"
	"bufio"
	"strings"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Please enter the port like \":<port>\"")
	var port string
	fmt.Scan("%s",&port)
	conn,err := net.Dial("tcp",port)
	if err != nil {
		log.Fatal("Connection to the server failed:",err,"\tport:",port)
	}


	defer conn.Close()
	inputReader := bufio.NewReader(os.Stdin)
	for {
		input,_ := inputReader.ReadString("\n")
		inputInfo := strings.Trim(input,"\n")//linux中使用`\n`，而Windows中是用`\r\n`
		if strings.ToUpper(inputInfo) == "Q" {
			return
		}//一定要事先设置好什么时候会退出程序
	}

	——,err := conn.Write([]byte(inputInfo))
	if err != nil {
		log.Fatal("write failed",err)
	}

	buf := [200000]byte{}
	n,err := conn.Read(buf[:])

	if err != nil {
		fmt.Println("recv failed,err:",err)
		log.Fatal("recv failed:",err)
	}

	fmt.Println(string(buf[:n]))
}