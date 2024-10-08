package process

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
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

func (header *dnsHeader) Flag(QR uint16, Opcode uint16, AA uint16, TC uint16, RD uint16, RA uint16, Rcode uint16) {
	header.Bits = QR<<15 + Opcode<<11 + AA<<10 + TC<<9 + RD<<8 + RA<<7 + Rcode
}

// 构建问题部分（域名）
type dnsQuery struct {
	Qutype  uint16
	Quclass uint16
}

func ParseDN(domain string) []byte { //这里的byte是单数！表示一个整体
	var (
		buffer   bytes.Buffer
		segments []string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		//第一句话表示将长度写入，第二局话表示将域名写入
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))

	return buffer.Bytes()
}

func (u *Udp) Send(dnsServer, domain string) (bytes.Buffer, uint16, int, *net.UDPConn, time.Duration, error) {
	var (
		conn       net.Conn
		err        error
		buffer     bytes.Buffer
		typechoice string
		RDchoice   int
	)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var randomId uint16 = uint16(rng.Intn(32768))

	requestHeader := dnsHeader{
		Id:      randomId,
		Qucount: 1,
		Ancount: 0,
		Nscount: 0,
		Adcount: 0,
	}
	fmt.Println("Do you desire to recurse request?(1.true/0.false)")
	_, err = fmt.Scanf("%d", &RDchoice)
	if err != nil || (RDchoice != 1 && RDchoice != 0) {
		log.Println("Unexpected content.We will open the desired recursion by default.")
		RDchoice = 1
	}
	requestHeader.Flag(0, 0, 0, 0, uint16(RDchoice), 0, 0)

	fmt.Println("Which type of the record do you want to check?(A,AAAA,NS,CNAME,MX,TXT)")
	fmt.Scanln(&typechoice)
	var requestQuery dnsQuery
	switch strings.ToUpper(strings.TrimSpace(typechoice)) {
	case "A":
		requestQuery = dnsQuery{
			Qutype:  1,
			Quclass: 1,
		}
	case "NS":
		requestQuery = dnsQuery{
			Qutype:  2,
			Quclass: 1,
		}
	case "CNAME":
		requestQuery = dnsQuery{
			Qutype:  5,
			Quclass: 1,
		}
	case "MX":
		requestQuery = dnsQuery{
			Qutype:  15,
			Quclass: 1,
		}
	case "TXT":
		requestQuery = dnsQuery{
			Qutype:  16,
			Quclass: 1,
		}
	case "AAAA":
		requestQuery = dnsQuery{
			Qutype:  28,
			Quclass: 1,
		}
	default:
		fmt.Println("We don't support the type of question.Instead we wil ask for A record for you.")
		requestQuery = dnsQuery{
			Qutype:  1,
			Quclass: 1,
		}
	}
	//请求的域名
	conn, err = net.Dial("udp", dnsServer)
	realconn, _ := conn.(*net.UDPConn)

	if err != nil {
		log.Printf("Failed to connect:%v\n", err.Error())
		return bytes.Buffer{}, 0, 0, realconn, time.Duration(0), err
	}

	binary.Write(&buffer, binary.BigEndian, requestHeader)
	binary.Write(&buffer, binary.BigEndian, ParseDN(domain))
	binary.Write(&buffer, binary.BigEndian, requestQuery)

	t1 := time.Now()
	_, err = conn.Write(buffer.Bytes())

	if err != nil {
		log.Printf("Failed to send the DNS query:%v\n", err.Error())
		return bytes.Buffer{}, 0, 0, realconn, time.Duration(0), err
	}
	requestLength := buffer.Len()

	duration := time.Since(t1)

	return buffer, randomId, requestLength, realconn, duration, nil
}

func (t *Tcp) Send(dnsServer, domain string) (bytes.Buffer, uint16, int, net.Conn, time.Duration, error) {
	var (
		conn       net.Conn
		err        error
		buffer     bytes.Buffer
		typechoice string
		RDchoice   int
	)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	var randomId uint16 = uint16(rng.Intn(32768))

	requestHeader := dnsHeader{
		Id:      randomId,
		Qucount: 1,
		Ancount: 0,
		Nscount: 0,
		Adcount: 0,
	}
	fmt.Println("Do you desire to recurse request?(1.true/0.false)")
	_, err = fmt.Scanf("%d", &RDchoice)
	if err != nil || (RDchoice != 1 && RDchoice != 0) {
		log.Println("Unexpected content.We will open the desired recursion by default.")
		RDchoice = 1
	}
	requestHeader.Flag(0, 0, 0, 0, uint16(RDchoice), 0, 0)

	fmt.Println("Which type of the record do you want to check?(A,AAAA,NS,CNAME,MX,TXT)")
	fmt.Scanln(&typechoice)
	var requestQuery dnsQuery
	switch strings.ToUpper(strings.TrimSpace(typechoice)) {
	case "A":
		requestQuery = dnsQuery{
			Qutype:  1,
			Quclass: 1,
		}
	case "NS":
		requestQuery = dnsQuery{
			Qutype:  2,
			Quclass: 1,
		}
	case "CNAME":
		requestQuery = dnsQuery{
			Qutype:  5,
			Quclass: 1,
		}
	case "MX":
		requestQuery = dnsQuery{
			Qutype:  15,
			Quclass: 1,
		}
	case "TXT":
		requestQuery = dnsQuery{
			Qutype:  16,
			Quclass: 1,
		}
	case "AAAA":
		requestQuery = dnsQuery{
			Qutype:  28,
			Quclass: 1,
		}
	default:
		fmt.Println("We don't support the type of question.Instead we wil ask for A record for you.")
		requestQuery = dnsQuery{
			Qutype:  1,
			Quclass: 1,
		}
	}
	//请求的域名
	conn, err = net.Dial("udp", dnsServer)
	realconn, _ := conn.(*net.UDPConn)

	if err != nil {
		log.Printf("Failed to connect:%v\n", err.Error())
		return bytes.Buffer{}, 0, 0, realconn, time.Duration(0), err
	}
	/*
		if dnsServer == "192.168.56.1:53" {
			binary.Write(&buffer, binary.BigEndian, uint8(2))
		}
	*/
	binary.Write(&buffer, binary.BigEndian, requestHeader)
	binary.Write(&buffer, binary.BigEndian, ParseDN(domain))
	binary.Write(&buffer, binary.BigEndian, requestQuery)

	t1 := time.Now()
	_, err = conn.Write(buffer.Bytes())

	if err != nil {
		log.Printf("Failed to send the DNS query:%v\n", err.Error())
		return bytes.Buffer{}, 0, 0, realconn, time.Duration(0), err
	}
	requestLength := buffer.Len()

	duration := time.Since(t1)

	return buffer, randomId, requestLength, realconn, duration, nil
}
