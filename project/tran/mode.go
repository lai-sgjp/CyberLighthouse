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
// bufio.Reader用完一定要重置！
func (t *Tcp) textMode(conn interface{}, message string) error { //这里的conn就相当于创建了实例，后面赋给了realconn
	realconn, _ := conn.(net.Conn)
	//message := string(buf[:])
	fmt.Println(message)

	_, err := realconn.Write([]byte("Successfully received!From the server.\n")) //注意这里将字符串转为[]byte
	if err != nil {
		realconn.Close()
		return err
	}
	return nil
}

func (t *Tcp) fileMode(conn interface{}, scanner *bufio.Reader) { //是不是因为是两个不同的reader导致？
	realconn, ok := conn.(net.Conn)
	defer scanner.Reset(scanner)
	//scanner := bufio.NewReader(realconn)
	if ok {
		//读取文件名
		filename, err := scanner.ReadString('\n')
		log.Println("init_filename######", filename)
		if err != nil {
			log.Println("Failed to get the file name.Please try again:", err.Error())
			return
		}
		filename = filename[:len(filename)-1] //注意这里有一个换行符是不需要的
		//读取文件大小
		var fileSize int64
		binary.Read(realconn, binary.BigEndian, &fileSize) //这样就可以直接读数字了
		log.Println("filesize######", fileSize)

		file, err := os.Create("received_" + filename)
		if err != nil {
			log.Printf("Failed to create the file:%v\n", err)

		}
		defer file.Close()

		//复制文件内容
		_, err = io.CopyN(file, realconn, fileSize)
		if err != nil {
			log.Printf("Failed to read the infomation:%v\n", err)
			//realconn.Close() //记得关闭连接
			return
		}

		fmt.Printf("Received the file \"%s\" Successfully\n", filename)
	}
	/*	_, err := realconn.Write([]byte("Successfully received!From the server.\n")) //注意这里将字符串转为[]byte
		if err != nil {
			realconn.Close()
			return
		}*/
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
	/*
		_, err = realconn.Write([]byte("Successfully received!From the server.\n")) //注意这里将字符串转为[]byte
		if err != nil {
			realconn.Close()
			return err
		}*/
	return nil
}

func (u *Udp) fileMode(conn interface{}) {
	realconn, ok := conn.(*net.UDPConn)
	scanner := bufio.NewReader(realconn)
	defer scanner.Reset(scanner)
	if ok {
		//读取文件名
		filename, err := scanner.ReadString('\n')
		log.Println("init_filename######", filename)
		if err != nil {
			log.Println("Failed to get the file name.Please try again:", err.Error())
			return
		}
		filename = filename[:len(filename)-1] //注意这里有一个换行符是不需要的
		//读取文件大小
		var fileSize int64
		binary.Read(realconn, binary.LittleEndian, &fileSize)
		log.Println("filesize######", fileSize)

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
			//realconn.Close() //记得关闭连接
			return
		}
		fmt.Printf("Received the file \"%s\" Successfully\n", filename)
		/*
			fmt.Printf("Received the file \"%s\" Successfully\n", filename)
				_, err = realconn.Write([]byte("Successfully received!From the server.\n")) //注意这里将字符串转为[]byte
				if err != nil {
					realconn.Close()
					return
				}*/
	}

}
