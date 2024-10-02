package tran_c
import (
	"bufio"
	"io"
	"os"
	"strings"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type quitError struct {
	message string
}
func (e quitError) Error() string{
	return e.message
}
//要发送什么（发送的模式）
func message() (string, error) {
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


//默认发送文本或者整数模式
func defaultSend(conn net.Conn)error {
	typeinfo := fmt.Sprintf("%d",0)
	_,err := conn.Write([]byte(typeinfo))
	if err != nil {
		log.Println(err)
		return err
	}

	trimmedInput, err := message()
	if err != nil {
		log.Println("Fail to input the infomation:", err,"\n")
		return err
	}
	//判断客户端退出条件
	if strings.ToUpper(trimmedInput) == "Q" {
		conn.Close() //记得关闭连接

		return quitError{message:"client has quitted"}
	}
	//发送数据
	_, err = conn.Write([]byte(trimmedInput))
	if err != nil {
		log.Println("Failed to send the message:",err,"\n")
	}
	return nil
}
//发送文件形式
//先发送文件名，再发送文件大小，最后发送文件内容
func SendFile(conn net.Conn,filename string) {
	//先发送一个数字看看是否是发送文件(编号为1)（协议请求头就是先给一个编码来判断是否哪种类型）
	typeinfo := fmt.Sprintf("%d",1)
	_,err := conn.Write([]byte(typeinfo))
	if err != nil {
		log.Fatal(err)
	}

	file,err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()//这里创建通道的时候又重复进行了选择，感觉不太高效。留一个坑后面有时间再填

	fmt.Fprintf(conn,filename+"\n")//向通道指定输出
	filestatus,err := file.Stat()
	if err != nil {
		log.Fatal("File is abnormal:",err)
	}
	fileSize := filestatus.Size()
	binary.Write(conn,binary.LittleEndian,fileSize)//小端法传输（小端在前）

	_,err = io.Copy(conn,file)
	if err != nil {
		log.Fatal("Failed to send the file:",err)
	}
}


//创立指定连接TCP
func CreateTCPConn(addr string) net.Conn {
	conn ,err := net.Dial("tcp",addr)

	if err != nil {
		//fmt.Printf("Failed to dial:", err)
		log.Fatal("Failed to dial:", err)
	}
		
	
	return conn
}
//创立指定连接UDP
func CreateUDPConn(addr string) net.Conn {
	conn ,err := net.Dial("udp",addr)

	if err != nil {
		//fmt.Printf("Failed to dial:", err)//原来这是两个都会打印到屏幕上的
		log.Fatal("Failed to dial:", err)
	}
	return conn
}
//选择某一个指定的连接(默认模式下)
//文件模式在client.go中已经定义
func Choose(tran string,addr string,message string) {
/*http建立在TCP之上
	url := fmt.Sprintf("%s/message",addr)
	//创建请求体
	body := bytes.NewBuffer([]byte(message))
	//创建请求
	req,err := http.NewRequest("GET",url,body)
	if err !+ nil {
		log.Fatal("Failed to set up the new request:",err)
	}
	//创建请求头（可选）
	req.Header.Set("Content_Type","text/plain")
	//初始化客户端
	client := &http.Client{}
	resp,err := client.Do(req)

*/
	switch tran {
	case "tcp":
		conn := CreateTCPConn(addr)
		defer conn.Close()
		//message,err := sendMessage()
		//if err != nil {
		//	log.Fatal("Failed to read the message:",err)
		//}
		err := defaultSend(conn)
		if err != nil {
			break
		}
	case "udp":
		conn := CreateUDPConn(addr)
		defer conn.Close()
		//message,err := sendMessage()
		//if err != nil {
		//	log .Fatal("Failed to read the message:",err)
		//}
		err := defaultSend(conn)
		if err != nil {
			break
		}
	}
}


//发送信息(默认)因为是循环的，所以难以拆分（或者说没思路）
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
		
		conn := CreateTCPConn(addr)
		defer conn.Close() //记得关闭连接
		for {
			err := defaultSend(conn)
			if err != nil {
				break
			}
		}
		
	case 1:
		/*
		addr, err := net.ResolveUDPAddr("udp", ":"+port)
		if err != nil {
			log.Fatal("Fail to get the correct TCP address:", err)
		}
		*/
		addr := fmt.Sprintf("%s:%s",host,port)
		conn := CreateUDPConn(addr)
		
		defer conn.Close() //记得关闭连接

		for {
			err := defaultSend(conn)
			if err != nil {
				break
			}
		}
		
	default:
		//fmt.Println("Input Error!")
		log.Fatal("Input Error!")
	}
}