package tran

import (
	"fmt"
	"log"
	"net"
	"os"
	"encoding/binary"
	"bufio"
	"io"
)

//tcp模式下
//埋下一个坑：提示客户端发送成功
func textMode(conn net.Conn,newBuf []byte)error{
	n, err := conn.Read(newBuf)
	if err != nil {
		log.Printf("Failed to recieve message from the client:%v\n", err)
		conn.Close() 

	}
	message := string(newBuf[:n-1])
	fmt.Printf("%s\n", message)

	_, err = conn.Write(newBuf[:n])
	if err != nil {
		log.Printf("Failed to send message to the client:%v\n", err)
		conn.Close() 
		return err
	}
	return nil
}

func fileMode(conn net.Conn)  {
	scanner := bufio.NewReader(conn)
	//读取文件名
	filename,_ := scanner.ReadString('\n')
	filename = filename[:len(filename) - 1]//注意这里有一个换行符是不需要的
	//读取文件大小
	var fileSize int64
	binary.Read(conn,binary.LittleEndian,&fileSize)

	file,err := os.Create("received_" + filename)
	if err != nil {
		log.Printf("Failed to create the file:%v\n",err)

	}
	defer file.Close()

	//复制文件内容
	_,err = io.CopyN(file,conn,fileSize)
	if err != nil {
		log.Printf("Failed to read the infomation:%v\n",err)

	}

	fmt.Printf("Received the file \"%s\" Successfully\n",filename)
}

//udp模式下
//埋下一个坑：下客户端发送成功
func textModeudp(conn *net.UDPConn,newBuf []byte) error {
	n, clientAddr,err := conn.ReadFromUDP(newBuf)
	if err != nil {
		log.Printf("Failed to recieve message from the client:%v\n", err)
		conn.Close() 
		return err
	}
	message := string(newBuf[:n])
	fmt.Printf("%s", message)

	_, err = conn.WriteToUDP(newBuf[:n],clientAddr)
	if err != nil {
		log.Printf("Failed to send message to the client:%v\n", err)
		conn.Close() 
		return err
	}
	return nil
}

func fileModeudp(conn *net.UDPConn){
	scanner := bufio.NewReader(conn)
	//读取文件名
	filename,_ := scanner.ReadString('\n')
	filename = filename[:len(filename) - 1]//注意这里有一个换行符是不需要的
	//读取文件大小
	var fileSize int64
	binary.Read(conn,binary.LittleEndian,&fileSize)

	file,err := os.Create("received_" + filename)
	if err != nil {
		log.Printf("Failed to create the file:%v\n",err)
		conn.Close() //记得关闭连接
	}
	defer file.Close()

	//复制文件内容
	_,err = io.CopyN(file,conn,fileSize)
	if err != nil {
		log.Printf("Failed to read the infomation:%v\n",err)
		conn.Close() //记得关闭连接
	}

	fmt.Printf("Received the file \"%s\" Successfully\n",filename)
}