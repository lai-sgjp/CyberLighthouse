package tran_c

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// 要发送什么（发送的模式）
func Message() (string, error) {

	fmt.Println("What information do you want to send?Type Q to quit.")
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n')
	if err != nil {
		return "", err
	}
	trimmedInput := strings.Trim(input, "\n")
	return trimmedInput, nil
}

// 默认发送文本或者整数模式
func defaultSend(conn net.Conn, message string) error { //字符串本身就是切片

	//判断客户端退出条件
	if strings.ToUpper(message) == "Q" {
		conn.Close()

		return quitError{message: "client has quitted"}
	}
	//发送数据
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Failed to send the message:", err.Error())
	}
	return nil
}

func defaultSendUDP(conn *net.UDPConn, message string) error {
	if strings.ToUpper(message) == "Q" {
		conn.Close()

		return quitError{message: "client has quitted"}
	}
	//发送数据
	_, err := conn.Write([]byte(message))
	if err != nil {
		log.Println("Failed to send the message:", err.Error())
		return err
	}
	return nil
}

// 发送文件形式
// 先发送文件名，再发送文件大小，最后发送文件内容
func SendFile(conn net.Conn, filename string) {
	//先发送一个数字看看是否是发送文件(编号为1)（协议请求头就是先给一个编码来判断是否哪种类型）
	typeinfo := fmt.Sprintf("%d\n", 1)
	_, err := conn.Write([]byte(typeinfo))
	if err != nil {
		log.Fatal("Failed to send the type", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to send the filename", err)
	}
	defer file.Close() //这里创建通道的时候又重复进行了选择，感觉不太高效。留一个坑后面有时间再填

	fmt.Fprintf(conn, filename+"\n") //向通道指定输出
	filestatus, err := file.Stat()
	if err != nil {
		log.Fatal("File is abnormal:", err)
	}
	fileSize := filestatus.Size()
	binary.Write(conn, binary.BigEndian, fileSize)
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal("Failed to send the file:", err)
	}
}

func SendFileUDP(conn *net.UDPConn, filename string) {
	//先发送一个数字看看是否是发送文件(编号为1)（协议请求头就是先给一个编码来判断是否哪种类型）
	typeinfo := fmt.Sprintf("%d\n", 1)
	_, err := conn.Write([]byte(typeinfo))
	if err != nil {
		log.Fatal("Failed to send the type", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to send the filename", err)
	}
	defer file.Close() //这里创建通道的时候又重复进行了选择，感觉不太高效。留一个坑后面有时间再填

	fmt.Fprintf(conn, filename+"\n") //向通道指定输出
	filestatus, err := file.Stat()
	if err != nil {
		log.Fatal("File is abnormal:", err)
	}
	fileSize := filestatus.Size()
	binary.Write(conn, binary.BigEndian, fileSize)
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Fatal("Failed to send the file:", err)
	}
}
