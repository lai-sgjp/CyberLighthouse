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

// 自定义错误
type quitError struct {
	message string
}

func (e quitError) Error() string {
	return e.message
}

// 要发送什么（发送的模式）
func Message() (string, error) {

	fmt.Println("What information do you want to send?Type Q to quit.")
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n') //读到包括\n
	if err != nil {
		return "", err
	}
	trimmedInput := strings.Trim(input, "\n")
	return trimmedInput, nil
}

func getmessage(message string) (string, error) {
	var err error
	if message == "" {
		message, err = Message()
		if err != nil {
			return "", err
		}
	}
	return message, nil
}

// 默认发送文本或者整数模式
func (t *Tcp) defaultSend(conn interface{}, message string) error { //字符串本身就是切片
	realconn, _ := conn.(net.Conn)
	//判断客户端退出条件
	if strings.ToUpper(message) == "Q" {
		realconn.Close()
		return quitError{message: "client has quitted"}
	}
	//发送数据类型
	typeinfo := fmt.Sprintf("%d"+message+"\n", 0)
	_, err := realconn.Write([]byte(typeinfo))
	if err != nil {
		log.Println("Failed to send the type and the message", err.Error())
		return err
	}
	//发送数据
	_, err = realconn.Write([]byte(message))
	if err != nil {
		log.Println("Failed to send the message:", err.Error())
	}
	return nil
}

func (u *Udp) defaultSend(conn interface{}, message string) error {
	realconn, _ := conn.(*net.UDPConn)
	if strings.ToUpper(message) == "Q" {
		realconn.Close()

		return quitError{message: "client has quitted"}
	}
	//发送数据类型和内容
	typeinfo := fmt.Sprintf("%d"+message+"\n", 0)
	_, err := realconn.Write([]byte(typeinfo))
	if err != nil {
		log.Println("Failed to send the type and the message:", err.Error())
		return err
	}
	return nil
}

// 发送文件形式
// 先发送文件名，再发送文件大小，最后发送文件内容
func (t *Tcp) SendFile(conn interface{}, filename string) {
	realconn, _ := conn.(net.Conn)
	//先发送一个数字看看是否是发送文件(编号为1)（协议请求头就是先给一个编码来判断是否哪种类型）
	typeinfo := fmt.Sprintf("%d\n", 1)
	_, err := realconn.Write([]byte(typeinfo))
	if err != nil {
		log.Fatal("Failed to send the type", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to send the filename", err)
	}
	defer file.Close() //这里创建通道的时候又重复进行了选择，感觉不太高效。留一个坑后面有时间再填

	fmt.Fprintf(realconn, filename+"\n") //向通道指定输出
	filestatus, err := file.Stat()
	if err != nil {
		log.Fatal("File is abnormal:", err)
	}
	fileSize := filestatus.Size()
	binary.Write(realconn, binary.BigEndian, fileSize)
	_, err = io.Copy(realconn, file)
	if err != nil {
		log.Fatal("Failed to send the file:", err)
	}
}

func (u *Udp) SendFile(conn interface{}, filename string) {
	realconn, _ := conn.(*net.UDPConn)
	//先发送一个数字看看是否是发送文件(编号为1)（协议请求头就是先给一个编码来判断是否哪种类型）
	typeinfo := fmt.Sprintf("%d\n", 1)
	_, err := realconn.Write([]byte(typeinfo))
	if err != nil {
		log.Fatal("Failed to send the type", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to send the filename", err)
	}
	defer file.Close() //这里创建通道的时候又重复进行了选择，感觉不太高效。留一个坑后面有时间再填

	fmt.Fprintf(realconn, filename+"\n") //向通道指定输出
	filestatus, err := file.Stat()
	if err != nil {
		log.Fatal("File is abnormal:", err)
	}
	fileSize := filestatus.Size()
	binary.Write(realconn, binary.BigEndian, fileSize)
	_, err = io.Copy(realconn, file)
	if err != nil {
		log.Fatal("Failed to send the file:", err)
	}
}
