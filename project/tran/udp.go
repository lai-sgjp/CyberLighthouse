package tran

import (
	"fmt"
	"log"
	"net"
	"os"
	"encoding/binary"//注意前面有encoding
	"bufio"
	"io"

)

func CreateUDPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveUDPAddr("udp", ":"+port)
		if err != nil {
			log.Printf("Failed to get the correct UDP address %s: %v\n", port, err)
			continue
		}

		//开始监听(或者说是连接)
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Printf("Failed to listen to UDP port %s:%v\n", port, err)
			continue
		}

		go ProcessUDP(conn)
	}
}



func ProcessUDP(conn *net.UDPConn) { //是一个指针(UDP与TCP还是不一样的)
	//TCP则是强调传输的安全，怎么说就是用Accept确认一次并建立专门的通道/socket
	//UDP是只要传输包就可以了，所以只需要一个指向某个通道的指针、地址
	defer conn.Close()
	buf := make([]byte, 8)
	n,_,err := conn.ReadFromUDP(buf)
	typeinfo := string(buf[:n])
	if err != nil {
		log.Printf("Failed to recieve type from the client:%v\n", err)
		conn.Close() //记得关闭连接
		return
	}
	switch typeinfo {
	case "0":
		newBuf := make([]byte,2048)
		for {
			err := textModeudp(conn,newBuf)
			if err != nil {
				break
			}
		}
	case "1":
		for {
			fileModeudp(conn)
		}
	default:
		log.Printf("Disapproval type:%v\n",err)
		conn.Close() //记得关闭连接
	}

}

//如果启用了go mod，一个包里面的才能互相看到

//埋下一个坑：下客户端发送成功
func textModeudp(conn *net.UDPConn,newBuf []byte) error {//UDP中conn的类型会加上net.UDPConn,而且注意会是指针
	n, clientAddr,err := conn.ReadFromUDP(newBuf)//这里的意思是将conn里面的读取到buf中
	if err != nil {
		log.Printf("Failed to recieve message from the client:%v\n", err)
		conn.Close() //记得关闭连接
		return err
	}
	message := string(newBuf[:n])
	fmt.Printf("%s", message)

	_, err = conn.WriteToUDP(newBuf[:n],clientAddr)
	if err != nil {
		log.Printf("Failed to send message to the client:%v\n", err)
		conn.Close() //记得关闭连接
		return err
	}

	return nil//无异常记得要返回nil
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