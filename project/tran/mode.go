package tran

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// tcp模式下
// 埋下一个坑：提示客户端发送成功
func (t *Tcp) textMode(conn interface{}, message string) error { //这里的conn就相当于创建了实例，后面赋给了realconn
	realconn, _ := conn.(net.Conn)
	//message := string(buf[:])
	fmt.Printf("%s\n", message)

	_, err := realconn.Write([]byte("Successfully received!From the server.")) //注意这里将字符串转为[]byte
	if err != nil {
		realconn.Close()
		return err
	}
	return nil
}

func (t *Tcp) fileMode(conn interface{}) {
	realconn, _ := conn.(net.Conn)
	scanner := bufio.NewReader(realconn)
	//读取文件名
	filename, _ := scanner.ReadString('\n')
	filename = filename[:len(filename)-1] //注意这里有一个换行符是不需要的
	//读取文件大小
	var fileSize int64
	binary.Read(realconn, binary.BigEndian, &fileSize)

	file, err := os.Create("received_" + filename)
	if err != nil {
		log.Printf("Failed to create the file:%v\n", err)

	}
	defer file.Close()

	//复制文件内容
	_, err = io.CopyN(file, realconn, fileSize)
	if err != nil {
		log.Printf("Failed to read the infomation:%v\n", err)

	}

	fmt.Printf("Received the file \"%s\" Successfully\n", filename)
}

// udp模式下
// 埋下一个坑：下客户端发送成功
func (u *Udp) textMode(clientAddr *net.UDPAddr, conn interface{}, message string) error {
	realconn, _ := conn.(*net.UDPConn)
	fmt.Println(message)
	log.Println("seccess read")

	_, err := realconn.WriteToUDP([]byte(message), clientAddr)
	if err != nil {
		realconn.Close()
		return err
	}
	return nil
}

func (u *Udp) fileMode(conn interface{}) {
	realconn, _ := conn.(*net.UDPConn)
	scanner := bufio.NewReader(realconn)
	//读取文件名
	filename, _ := scanner.ReadString('\n')
	filename = filename[:len(filename)-1] //注意这里有一个换行符是不需要的
	//读取文件大小
	var fileSize int64
	binary.Read(realconn, binary.LittleEndian, &fileSize)

	file, err := os.Create("received_" + filename)
	if err != nil {
		log.Printf("Failed to create the file:%v\n", err)
		realconn.Close() //记得关闭连接
	}
	defer file.Close()

	//复制文件内容
	_, err = io.CopyN(file, realconn, fileSize)
	if err != nil {
		log.Printf("Failed to read the infomation:%v\n", err)
		realconn.Close() //记得关闭连接
	}

	fmt.Printf("Received the file \"%s\" Successfully\n", filename)
}
