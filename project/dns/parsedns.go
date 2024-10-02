package main

import (
	"bufio"
	"net"
	"bytes"
	"encoding/binary"
	"strings"
	"log"
)

func parase(responseLength,sendId int,buf bytes.Buffer,conn net.Conn) (map[string]string,error) {
	var answerbuf []bytes
	answerLength,err := conn.Read(answerbuf)

	//对响应报文进行检查
	if err != nil {
		log.Printf("Failed to receive the response from the DNS server:%v\n",err.Error())
		return
	}
	recvId := answerbuf[:2]
	if recvId != sendId || answerbuf <= responseLength{
		log.Printf("Failed to receive the correct response from the DNS server:%v\n",err.Error())
	}

	//开始解析
	i := 2
	
	//解析头部
	//注意：i++是不能直接放在其他语句块的里面（这与C很不同）
	QR := answerbuf[i:i+1]
	i++
	Opcode := answerbuf[i:i+2]
	i += 2

}
