package process

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func caseJustify(a uint16) bool {
	switch a {
	case 0:
		return false
	case 1:
		return true
	}
	return false
}

type DnsHeader struct {
	RecvId                             uint16
	AA, TC, RD, RA                     bool
	Rcode, Opcode                      string
	Qucount, Aucount, Ancount, Adcount uint16
}
type Question struct {
	Domain                      string
	Questionclass, Questiontype uint16
}
type Resource struct {
	TTL  uint32
	Data string
}
type Returninfomation struct {
	Other       bytes.Buffer
	Err         error
	DNSHeader   DnsHeader
	DNSQuestion Question
	DNSResource []Resource
}

// 如果使用空接口，必须先进行类型断言在进行数值接收
func (t *Tcp) Parse(queryLength int, sendId uint16, conn interface{}) chan Returninfomation {
	realconn, _ := conn.(net.Conn)
	t1 := time.Now()

	ch := make(chan Returninfomation, 1)

	answerbuf := make([]byte, 512)
	answerLength, err := realconn.Read(answerbuf)

	//对响应报文进行检查
	if err != nil {
		log.Printf("Failed to receive the response from the DNS server:%v\n", err.Error())
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}

	Ancount, header, err := parseHeader(sendId, queryLength, answerLength, answerbuf[:12])
	if err != nil {
		log.Printf("Failed to parse the header:%v\n", err.Error())
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}

	domain, choice, questionlength, question := parseQuestion(answerbuf[12:])
	//到底这里应该用uint16还是int?之后统一改成uint16吧（埋一个坑）
	resourcelength, ResourceList := parseResource(domain, choice, Ancount, questionlength, answerbuf[12:])
	if len(answerbuf) == 12+questionlength+resourcelength {
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}
	other, err := storeOther(answerbuf[12+questionlength+resourcelength:])
	if err != nil {
		log.Printf("Failed to store the Authority and Additional infomation")
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}
	returninfo := Returninfomation{
		Other:       other,
		Err:         nil,
		DNSHeader:   header,
		DNSQuestion: question,
		DNSResource: ResourceList,
	}
	ch <- returninfo
	close(ch)
	durataion := time.Since(t1)
	fmt.Println(durataion)
	return ch
}

func (u *Udp) Parse(queryLength int, sendId uint16, conn interface{}) chan Returninfomation {
	realconn, _ := conn.(*net.UDPConn)
	t1 := time.Now()

	ch := make(chan Returninfomation, 1)

	answerbuf := make([]byte, 512)
	answerLength, err := realconn.Read(answerbuf)

	//对响应报文进行检查
	if err != nil {
		log.Printf("Failed to receive the response from the DNS server:%v\n", err.Error())
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}

	Ancount, header, err := parseHeader(sendId, queryLength, answerLength, answerbuf[:12])
	if err != nil {
		log.Printf("Failed to parse the header:%v\n", err.Error())
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}

	domain, choice, questionlength, question := parseQuestion(answerbuf[12:])
	//到底这里应该用uint16还是int?之后统一改成uint16吧（埋一个坑）
	resourcelength, ResourceList := parseResource(domain, choice, Ancount, questionlength, answerbuf[12:])
	if len(answerbuf) == 12+questionlength+resourcelength {
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}
	other, err := storeOther(answerbuf[12+questionlength+resourcelength:])
	if err != nil {
		log.Printf("Failed to store the Authority and Additional infomation")
		returninfo := Returninfomation{
			Other: bytes.Buffer{},
			Err:   err,
		}
		ch <- returninfo
		close(ch)
		return ch
	}
	returninfo := Returninfomation{
		Other:       other,
		Err:         nil,
		DNSHeader:   header,
		DNSQuestion: question,
		DNSResource: ResourceList,
	}
	ch <- returninfo
	close(ch)
	durataion := time.Since(t1)
	fmt.Println(durataion)
	return ch
}

func parseHeader(sendId uint16, queryLength, answerLength int, answerbuf []byte) (uint16, DnsHeader, error) {
	recvId := uint16(answerbuf[0])<<8 + uint16(answerbuf[1]) //完全可以自动化不用手动
	if recvId != sendId || answerLength <= queryLength {
		log.Printf("Received ID:%v\tSend ID:%v", recvId, sendId)
		log.Printf("Failed to receive the correct response from the DNS server\n")
	}

	//开始解析

	//解析头部
	var AA, TC, RD, RA, Rcode, Opcode uint16
	//网络编程里面统一采用大端法
	//这里涉及一位一位，所以采用位运算
	flags := uint16(answerbuf[2])<<8 + uint16(answerbuf[3]) //因为是大端法，所以先读取rcode
	//这里原来可以直接用binary.Read(*byte.reader,encoding_way,buffer)啊
	Rcode = flags & 0xF
	flags >>= (4 + 3)
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
		return 0, DnsHeader{}, nil
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

	qucount := uint16(answerbuf[4])<<8 + uint16(answerbuf[5])
	ancount := uint16(answerbuf[6])<<8 + uint16(answerbuf[7])
	aucount := uint16(answerbuf[8])<<8 + uint16(answerbuf[9])
	adcount := uint16(answerbuf[10])<<8 + uint16(answerbuf[11])
	fmt.Println("Header:")
	fmt.Printf("status:%v\tid:%d\n", Rcode1, recvId)
	fmt.Printf("Opcode:%v\tauthoritative answer:%v\ttruncated :%v\nrecursion desired:%v\trecursion available:%v\n", Opcode1, AA1, TC1, RD1, RA1)
	fmt.Printf("question count:%d\tanswer count:%d\tauthority record count:%d\tadditional record count:%d\n", qucount, ancount, aucount, adcount)
	if RA == 0 {
		log.Println("This DNS server doesn't support the \"recursion\" function.")
	}
	header := DnsHeader{
		RecvId: recvId,
		AA:     AA1, TC: TC1, RD: RD1, RA: RA1,
		Rcode:   Rcode1,
		Opcode:  Opcode1,
		Qucount: qucount, Ancount: ancount, Aucount: aucount, Adcount: adcount,
	}
	return ancount, header, nil
}

func parseQuestion(answerbuf []byte) (string, string, int, Question) {
	var (
		domain                      string
		bits                        int
		i                           int
		questionclass, questiontype uint16
	)
	i = 0

	for {
		bits = int(answerbuf[i]) //表示长度的只会占1个字符！！
		i++
		if bits == 0 {
			break
		}
		readbuf := make([]byte, bits)
		scanner := bytes.NewReader(answerbuf[i : i+bits])
		err := binary.Read(scanner, binary.BigEndian, &readbuf)
		if err != nil {
			log.Println("Failed to read the domain:", err.Error())
		}
		domain += (string(readbuf) + ".")
		i += bits
	}

	questiontype = binary.BigEndian.Uint16(answerbuf[i : i+2])
	questionclass = binary.BigEndian.Uint16(answerbuf[i+2 : i+4])
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
	question := Question{
		Domain:        domain,
		Questiontype:  questiontype,
		Questionclass: questionclass,
	}
	return domain, questiontype1, i + 1, question //因为是长度，所以要在索引的基础上+1
}

func parseResource(domain, choice string, ancount uint16, questionlength int, answerbuf []byte) (int, []Resource) {
	fmt.Println("Answer Section:")
	var i int = 6 + questionlength //因为域名出现压缩，所以只用两个字节存储，后面的questiontype和questionclass一致，分别为2个字节
	var j uint16
	ttl := binary.BigEndian.Uint32(answerbuf[i : i+4]) //是4个字节！！！
	i += 4
	datalength := int(answerbuf[i])
	i++
	ResourceList := make([]Resource, ancount)
	for j = 0; j < ancount; j++ {
		var data string
		if choice == "A" {
			data = net.IPv4(answerbuf[i], answerbuf[i+1], answerbuf[i+2], answerbuf[i+3]).String()
			i += 4
		} else if choice == "AAAA" {
			for k := 0; k < 8; k++ {
				tmp := binary.BigEndian.Uint16(answerbuf[i : i+2])
				i += 2
				if k == 0 {
					data = fmt.Sprintf("%x", tmp)
					continue
				}
				data = fmt.Sprintf(data+":%x", tmp)
			}
		} else {
			tag := binary.BigEndian.Uint16(answerbuf[i : i+2])
			if (tag & 0xc0) == 0xc0 {
				i += 2
				continue
			}
			data = string(answerbuf[i : i+datalength])
			i += datalength
		}
		ResourceList[j] = Resource{
			TTL:  ttl,
			Data: data,
		}
		fmt.Printf("%s\t\t\tIN\t%s\t\t\t%v\t\t\tTTL:%vs\n", domain, choice, data, ttl)
	}
	return i, ResourceList
	//return i - (12 + questionlength), ResourceList
}

func storeOther(answerbuf []byte) (bytes.Buffer, error) {
	var other bytes.Buffer
	_, err := other.Write(answerbuf) //联想：二进制binary.Write(*otr->buffer,way,context)注意Write的W是大写的
	return other, err
}

//NS指针压缩是在不知道怎么搞，先留下来埋坑
/*
	var k int
	for k = 0; k < datalength; k++ {
		//判断是否进行了标签压缩
		tag := binary.BigEndian.Uint16(answerbuf[i+k : i+k+2])
		if (tag & 0xc0) == 0xc0 {
			bias := tag & 0x3f
			tmp := bias - 12
			for {
				bits := uint16(answerbuf[tmp]) //表示长度的只会占1个字符！！
				if bits == 0 {
					break
				}
				tmp++
				readbuf := make([]byte, bits)
				scanner := bytes.NewReader(answerbuf[tmp : tmp+bits])
				err := binary.Read(scanner, binary.BigEndian, &readbuf)
				if err != nil {
					log.Println("Failed to read the domain:", err.Error())
				}
				data += (string(readbuf) + ".")
				tmp += bits
			}
			data += (data + ".")
			k += 2
			continue
		}
		tmp := int(answerbuf[i+k])
		data += string(answerbuf[i+k:i+k+tmp]) + "."
		k += (1 + tmp)
	}
	i += datalength
	break
*/
