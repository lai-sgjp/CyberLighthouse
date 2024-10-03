package report

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	//"math/rand"
	"net"
	"strings"
	"time"
)

// 报文头部
type dnsHeader struct {
	Id                                 uint16
	Bits                               uint16
	Qucount, Ancount, Nscount, Adcount uint16
}

// 在函数名称前加上`(structure)`表示是结构体的方法
// 这里是为了实现flag的定制（因为相对于其他，flag比较灵活）
func (header *dnsHeader) Flag(QR uint16, Opcode uint16, AA uint16, TC uint16, RD uint16, RA uint16, Rcode uint16) {
	header.Bits = QR<<15 + Opcode<<11 + AA<<10 + TC<<9 + RD<<8 + RA<<7 + Rcode
} //这里的bits是标志的意思

// 构建问题部分（域名）
type dnsQuery struct {
	Qutype  uint16
	Quclass uint16
}

func ParseDN(domain string) []byte { //这里的byte是单数！表示一个整体
	var (
		buffer bytes.Buffer //是大写的B
		//bytes.Buffer是大小可变的字节管理区
		segments []string = strings.Split(domain, ".")
	)
	//这里一个var是一种编程风格
	//可以强调这个函数里面有哪些，生命周期就在这一个函数
	for _, seg := range segments {
		//第一句话表示将长度写入，第二局话表示将域名写入
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
		/*
					fileSize := filestatus.Size()
			binary.Write(conn,binary.LittleEndian,fileSize)
			_,err = io.Copy(conn,file)
		*/
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00)) //这里的0x00表示域名的结束by添加00的方式

	return buffer.Bytes()
}

func Send(dnsServer, domain string) (bytes.Buffer, uint16, int, net.Conn, time.Duration, error) { //最后一个表示一段时长，以秒作为单位

	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source) //提供种子
	var randomId uint16 = uint16(rng.Intn(32768))

	requestHeader := dnsHeader{
		Id:      randomId, //期待改成随机数
		Qucount: 1,
		Ancount: 0,
		Nscount: 0,
		Adcount: 0, //结构体赋值之后每一个后面都要加上逗号
	}

	requestHeader.Flag(0, 0, 0, 0, 1, 0, 0)

	//请求的域名
	requestQuery := dnsQuery{
		Qutype: 1, //这个表示A记录
		//另：28是IPV6地址
		//2为NS
		//5为重定向
		//15为电子邮件
		//16为文本信息
		Quclass: 1, //这个表示是TCP协议
	}

	var (
		conn   net.Conn
		err    error
		buffer bytes.Buffer
	)

	conn, err = net.Dial("udp", dnsServer) //后面传的是协议和字符串型地址

	if err != nil {
		log.Printf("Failed to connect:%v\n", err.Error())
		return bytes.Buffer{}, 0, 0, conn, time.Duration(0), err
	}
	//defer conn.Close() //defer是在它之下的代码执行之后有异常才会断，上面是没有异常的
	//这里不要关啊！后面还要收取的时候还要用的！

	//分3端写表示3个部分，连续追加。当然可以写在一串，这样写逻辑更加清晰
	//问题部分一定是name->type->class
	//注意，一定是大端法传输！
	binary.Write(&buffer, binary.BigEndian, requestHeader)
	binary.Write(&buffer, binary.BigEndian, ParseDN(domain))
	binary.Write(&buffer, binary.BigEndian, requestQuery)

	t1 := time.Now()
	_, err = conn.Write(buffer.Bytes())
	//buffer.Bytes是返回目前缓冲区所有的内容
	if err != nil {
		log.Printf("Failed to send the DNS query:%v\n", err.Error())
		return bytes.Buffer{}, 0, 0, conn, time.Duration(0), err
	}
	requestLength := buffer.Len() //为了解析而使用//注意这里是变长的，应该使用方法.Len()

	duration := time.Since(t1)
	//这里使用time.Since而不是用time.Now().Sub()
	/*
		responseLength, err:= conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read the DNS response:%v\n",err.Error())
			return make([]byte,0),0,0
		}
		duration := time.Now().Sub(t1)
	*/
	return buffer, randomId, requestLength, conn, duration, nil
}
