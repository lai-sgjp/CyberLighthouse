package main

import (
	"bufio"
	"net"
	"bytes"
	"encoding/binary"
	"strings"
	"log"
)

func parse(responseLength,sendId int,buf bytes.Buffer,conn net.Conn) error{
	answerLength,err := conn.Read(answerbuf)//这里的answerbuf是[]byte类型

	//对响应报文进行检查
	if err != nil {
		log.Printf("Failed to receive the response from the DNS server:%v\n",err.Error())
		return err
	}

	Ancount,err := parseHeader(answerbuf[:12])
	if err != nil {
		log.Printf("Failed to parse the header:%v\n",err.Error())
		return err
	}

	domain,questionlength := parseQuestion(answerbuf[12:])
	resourcelength := parseResource(domain,Ancount,answerbuf[12+questionlength:])
	if len(answerbuf) > 12+questionlength+resourcelengt {
		other := storeOther(answerbuf[12+questionlength+resourcelength:])
	}
}

func parseHeader(answerbuf [12]byte) (int,error){
	recvId := answerbuf[:2]
	if recvId != sendId || answerbuf <= responseLength{
		log.Printf("Failed to receive the correct response from the DNS server:%v\n",err.Error())
	}

	//开始解析
	
	//解析头部
	//注意：i++是不能直接放在其他语句块的里面（这与C很不同）
	var AA,TC,RD,RA,Rcode,Opcode interface{}

	flags := answerbuf[2:4]
	Rcode = flags & 0x1111
	flags >>= (4 + 3)//有保留位3位
	RA = flags & 0x1
	flags >>= 1
	RD = flags & 0x1
	flags >>= 1
	TC = flags & 0x1
	flags >>= 1
	AA = flags & 0x1
	flags >>= 1
	Opcode = flags & 0x1111
	flags >>= 4
	QR := flags
	func caseJustify(a interface{}) bool {
		switch a {
		case 0:
			return false
		case 1:
			return true
		}
	}
	if QR != 1 {
		log.Printf("Answer received is not the response from the DNS server:%v\n",err.Error())
	}
	switch Opcode {
	case 0:
		Opcode = "Standard query"
	case 1:
		Opcode = "reverse query"
	case 2:
		Opcode = "server status"
	}
	AA = caseJustify(AA)
	TC = caseJustify(TC)
	RD = caseJustify(RD)
	RA = caseJustify(RA)
	switch Rcode {
	case 0:
		Rcode = "No Error"
	case 2:
		Rcode = "Server fail"
	case 3:
		Rcode = "Name Error"
	}

	Qucount := answerbuf[4:6]
	Ancount := answerbuf[6:8]
	Aucount := answerbuf[8:10]
	Adcount := answerbuf[10:12]
	fmt.Println("Header:")
	fmt.Println("status:%v\tid:%d",Rcode,recvId)
	fmt.Println("Opcode:%v\tauthoritative answer:%v\ttruncated :%v\nrecursion desired:%v\trecursion available:%v",Opcode,AA,TC,RD,RA)
	fmt.Println("question count:%d\tanswer count:%d\tauthority record count:%d\tadditional record count:%d",Qucount,Ancount,Aucount,Adcount)
	return Ancount,nil
}
//暂且只支持查询一个
//func parseQuestion(answerbuf []byte,questiontype int,questionClass int) error {
func parseQuestion(answerbuf []byte) (string,int) {
	var(
		domain string//字符串可以通过`+`来连接
		bit int
		i int
		questionclass,questiontype interface{}//通过使用接口达到任意类型
	) 
	i = 0
	for {
		bits = answerbuf[:1]
		i++
		if bit == 0 {
			break
		}
		domain += (string(answerbuf[ i :i + bits]) + ".")
		i += bits//以ACSII码进行存储，一个英文字符占一个字节
	}
	questiontype = answerbuf[i:i+2]
	questionclass = answerbuf[i+2:i+4]//这两个都是2个字节
	i += 4
	switch questiontype {
	case 1:
		questiontype = "A"
	case 2:
		questiontype = "NS"
	case 5:
		questiontype = "CNAME"
	case 15:
		questiontype = "MX"
	case 16:
		questiontype = "TXT"
	case 28:
		questiontype = "AAAA"
	}
	if questionclass == 1 {
		questionclass = "TCP/IP"
	}

	fmt.Println("Quextion section:")
	fmt.Println("%s\t\t\tIN\t%v",domain,questiontype)
	fmt.Println("%v",questionclass)
	return domain,i+1
}
//暂且只支持查询A记录
func parseResource(domain string,Ancount int,answerbuf []byte) int {
	fmt.Println("Answer Section:")
	var i int = 6//因为域名出现压缩，所以只用两个字节存储，后面的questiontype和questionclass一致，分别为2个字节
	for j := 0;j < Ancount:j++ {
		TTL := answerbuf[i:i+4]
		datalength := answerbuf[i+4:i+6]
		i += 6
		data := string(answerbuf[i:i+datalength])
		i += (datalength + 6)
		fmt.Println("%s\t\t\tIN\tA\t\t\t%v",domain,data)
	}
}

func storeOther(answerbuf []byte)bytes.Buffer {
	var other bytes.Buffer
	_, err=other.write(answerbuf)//联想：二进制binary.Write(*otr->buffer,way,context)注意Write的W是大写的
	//这里的write相当于类的实现方法
	return other
}