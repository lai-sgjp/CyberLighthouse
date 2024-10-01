package main

import (
	"http/net"
	"fmt"
	"os"
	"bufio"
)


//这是通过net包实现的
func process(conn net.Conn) {
	defer conn.Close()//提前设置关闭连接，结束后执行

	for {
		reader := bufio.NewReader(conn)//这是一个初始化了的读取器
		var buf [128]byte//最多容纳128个英文字符
		n , err := reader.Read(buf[:])//读取网址数据

		if err != nil {
			fmt.Println("read from client failed,err:" ,err)
			break
		}

		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：",recvStr)
		conn.Write([]byte(recvStr))
	}
}

func main() {
	//建立监听通道
	listen,err := net.Listen("tcp","127.0.0.1:8080")
	if err != nil {
		fmt . Println("listen failed,err :", err)
		return
	}

	for {
		//伺机找到连接
		conn,err := listen.Accept()//建立连接
		if err != nil {
			fmt.Println("accept failed,err:",err)
			continue
		}

		go pocess(conn)
	}
}

//建立监听器`listen`->不停循环寻找连接->`accept`建立连接 ->{读取数据->格式化数据->发送数据`conn.Write`->读取数据`conn.Read`->打印数据 ->关闭连接}