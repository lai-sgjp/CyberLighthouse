package tran
//被引用的包名必须和文件名一致
import (
	"fmt"
	"log"
	"net"
	"os"
	"encoding/binary"//注意前面有encoding
	"bufio"
	"io"
)

func CreateTCPSer(ports []string) {
	for _, port := range ports {
		//建立端口地址
		addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)//addr是*net.TCPAddr类型
		if err != nil {
			log.Printf("Failed to get the correct TCP address %s: %v\n", port, err)
			continue
		}
		//监听端口地址(采用多协程)注意这里是创建一个而不是遍历所有
		listener, err := net.ListenTCP("tcp", addr)//这里的listener是*net.TCPListener数据类型
		if err != nil {
			log.Printf("Failed to listen to the TCP address %s:%v\n", port, err)
			continue
		}//建立监听不放进协程，因为建立可能失败，并且会有协程竞争，降低效率
		go func(listener *net.TCPListener) {
			var conn net.Conn
			var err error
			for {
				conn, err = listener.Accept()//这里的conn是net.Conn接口值，实际可能是`*net.TCPConn``net.UDPConn`
				if err != nil {
					log.Printf("Failed to accept the connection:%v\n", err)
					continue
				}
			}
			//没有错误表示正常监听（即监听已经建立）
			go ProcessTCP(conn) //使用go采用多协程的方式从而并发处理
		}(listener)//定义在外部又不是全局变量，匿名函数想引用必须进行捕获
		
	}
}

func ProcessTCP(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 8)

	n,err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to recieve type from the client:%v\n", err)
		conn.Close() //记得关闭连接
		return
	}

	typeinfo := string(buf[:n])
	switch typeinfo {
	case "0":
		newBuf := make([]byte,2048)
		for {
			err := textMode(conn,newBuf)
			if err != nil {
				break
			}
		}//遇到这种需要通过循环的函数里面错误来跳出循环的应该设置错误值来跳出
	case "1":
		fileMode(conn)
	default:
		log.Printf("Disapproval type:%v\n",err)
		conn.Close() //记得关闭连接
	}
	
}

//埋下一个坑：提示客户端发送成功
func textMode(conn net.Conn,newBuf []byte)error{
	n, err := conn.Read(newBuf)//这里的意思是将conn里面的读取到buf中
	if err != nil {
		log.Printf("Failed to recieve message from the client:%v\n", err)
		conn.Close() //记得关闭连接
		return err

	}
	message := string(newBuf[:n-1])
	fmt.Printf("%s\n", message)

	_, err = conn.Write(newBuf[:n])
	if err != nil {
		log.Printf("Failed to send message to the client:%v\n", err)
		conn.Close() //记得关闭连接
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