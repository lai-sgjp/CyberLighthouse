package report

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func caseJustify(a uint16) bool {
	switch a {
	case 0:
		return false
	case 1:
		return true
	}
	return false
} //接口判断必须先进行类型断言
//除了main函数外，不能再一个函数里面定义一个函数（最好main也不要）

func Parse(queryLength int, sendId uint16, conn net.Conn) (bytes.Buffer, error) {
	answerbuf := make([]byte, 512)
	answerLength, err := conn.Read(answerbuf) //这里的answerbuf是[]byte类型
	log.Println(answerbuf)

	//对响应报文进行检查
	if err != nil {
		log.Printf("Failed to receive the response from the DNS server:%v\n", err.Error())
		return bytes.Buffer{}, err //注意返回一个空这里是这样写的（依据类型的定义）
	}

	Ancount, err := parseHeader(sendId, queryLength, answerLength, answerbuf[:12])
	if err != nil {
		log.Printf("Failed to parse the header:%v\n", err.Error())
		return bytes.Buffer{}, err
	}

	domain, questionlength := parseQuestion(answerbuf[12:])
	//到底这里应该用uint16还是int?之后统一改成uint16吧（埋一个坑）
	resourcelength := parseResource(domain, Ancount, answerbuf[12+questionlength:])
	if len(answerbuf) == 12+questionlength+resourcelength {
		return bytes.Buffer{}, nil
	}
	other, err := storeOther(answerbuf[12+questionlength+resourcelength:])
	if err != nil {
		log.Printf("Failed to store the Authority and Additional infomation")
		return other, err
	}
	return bytes.Buffer{}, nil
}

func parseHeader(sendId uint16, queryLength, answerLength int, answerbuf []byte) (uint16, error) {
	recvId := uint16(answerbuf[0])<<8 + uint16(answerbuf[1]) //进行“移位”就是字节转为数字的方法
	if recvId != sendId || answerLength <= queryLength {     //将数字转为字节再进行比较||将字节转为数字
		log.Printf("Received ID:%v\nSend ID:%v", recvId, sendId)
		log.Printf("Failed to receive the correct response from the DNS server\n")
	}

	//开始解析

	//解析头部
	//注意：i++是不能直接放在其他语句块的里面（这与C很不同）
	var AA, TC, RD, RA, Rcode, Opcode uint16
	//网络编程里面统一采用大端法
	flags := uint16(answerbuf[2])<<8 + uint16(answerbuf[3]) //因为是大端法，所以先读取rcode
	//这里原来可以直接用binary.Read(*byte.reader,encoding_way,buffer)啊
	Rcode = flags & 0xF //0x表示16进制，所以不是0x1111，而是0xF(曲后4位)
	flags >>= (4 + 3)   //有保留位3位
	RA = flags & 0x1
	flags >>= 1
	RD = flags & 0x1
	flags >>= 1
	TC = flags & 0x1
	flags >>= 1
	AA = flags & 0x1
	flags >>= 1
	Opcode = flags & 0xF
	flags >>= 4
	QR := flags

	if QR != 1 {
		log.Printf("Answer received is not the response from the DNS server\n")
		return 0, nil
	}
	var Opcode1, Rcode1 string
	switch Opcode {
	case 0:
		Opcode1 = "Standard query"
	case 1:
		Opcode1 = "reverse query"
	case 2:
		Opcode1 = "server status"
	}
	AA1 := caseJustify(AA)
	TC1 := caseJustify(TC)
	RD1 := caseJustify(RD)
	RA1 := caseJustify(RA)
	switch Rcode {
	case 0:
		Rcode1 = "No Error"
	case 2:
		Rcode1 = "Server fail"
	case 3:
		Rcode1 = "Name Error"
	}

	Qucount := uint16(answerbuf[4])<<8 + uint16(answerbuf[5])
	Ancount := uint16(answerbuf[6])<<8 + uint16(answerbuf[7])
	Aucount := uint16(answerbuf[8])<<8 + uint16(answerbuf[9])
	Adcount := uint16(answerbuf[10])<<8 + uint16(answerbuf[11])
	fmt.Println("Header:")
	fmt.Printf("status:%v\tid:%d\n", Rcode1, recvId)
	fmt.Printf("Opcode:%v\tauthoritative answer:%v\ttruncated :%v\nrecursion desired:%v\trecursion available:%v\n", Opcode1, AA1, TC1, RD1, RA1)
	fmt.Printf("question count:%d\tanswer count:%d\tauthority record count:%d\tadditional record count:%d\n", Qucount, Ancount, Aucount, Adcount)
	return Ancount, nil
}

// 暂且只支持查询一个
// func parseQuestion(answerbuf []byte,questiontype int,questionClass int) error {
func parseQuestion(answerbuf []byte) (string, int) {
	var (
		domain                      string //字符串可以通过`+`来连接
		bit                         int
		i                           int
		questionclass, questiontype uint16 //通过使用接口达到任意类型
	)
	i = 0

	for {
		bits := int(answerbuf[0])
		i++
		if bit == 0 {
			break
		}
		var readbuf []byte
		scanner := bytes.NewReader(answerbuf[i : i+bits])       //这里创建的就是一个指针
		err := binary.Read(scanner, binary.BigEndian, &readbuf) //将二进制数据流转为字符串切片
		if err != nil {
			log.Println("Failed to read the domain:", err.Error())
		}
		domain += (string(readbuf) + ".")
		i += bits //以ACSII码进行存储，一个英文字符占一个字节
	}
	domain += "."

	questiontype = binary.LittleEndian.Uint16(answerbuf[i : i+2])
	questionclass = binary.LittleEndian.Uint16(answerbuf[i+2 : i+4]) //这两个都是2个字节
	i += 4
	var questiontype1, questionclass1 string
	switch questiontype {
	case 1:
		questiontype1 = "A"
	case 2:
		questiontype1 = "NS"
	case 5:
		questiontype1 = "CNAME"
	case 15:
		questiontype1 = "MX"
	case 16:
		questiontype1 = "TXT"
	case 28:
		questiontype1 = "AAAA"
	}
	if questionclass == 1 {
		questionclass1 = "TCP/IP"
	}

	fmt.Println("Question section:")
	fmt.Printf("%s\t\t\tIN\t%v\n", domain, questiontype1)
	fmt.Printf("%v\n", questionclass1)
	return domain, i + 1
}

// 暂且只支持查询A记录
func parseResource(domain string, Ancount uint16, answerbuf []byte) int {
	fmt.Println("Answer Section:")
	var i uint16 = 6 //因为域名出现压缩，所以只用两个字节存储，后面的questiontype和questionclass一致，分别为2个字节
	var j uint16
	for j = 0; j < Ancount; j++ {
		/*
			TTL := uint32(answerbuf[i])<<24 + uint32(answerbuf[i+1])<<16 + uint32(answerbuf[i+2])<<8 + uint32(answerbuf[i+3])
			var TTLbuf [4]byte
			scanner := binary.NewReader(answerbuf[i : i+5])
		*/
		scanner := bytes.NewReader(answerbuf[i : i+5])
		var TTLbuf [4]byte
		err := binary.Read(scanner, binary.BigEndian, &TTLbuf)
		if err != nil {
			log.Println("Failed to parse the TTL:", err)
		}
		TTL := string(TTLbuf[:4]) //不能直接将固定的字符组转为字符串，因为字符串本质是切片

		datalength := uint16(answerbuf[i+4])<<8 + uint16(answerbuf[i+5])
		i += 6
		data := string(answerbuf[i : i+datalength])
		i += (datalength + 6)
		fmt.Printf("%s\t\t\tIN\tA\t\t\t%v\t\t\tTTL:%vs\n", domain, data, TTL)
	}
	return int(i)
}

func storeOther(answerbuf []byte) (bytes.Buffer, error) {
	var other bytes.Buffer
	_, err := other.Write(answerbuf) //联想：二进制binary.Write(*otr->buffer,way,context)注意Write的W是大写的
	//这里的write相当于类的实现方法
	return other, err
}
