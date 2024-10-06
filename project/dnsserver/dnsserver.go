package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func clientConn()error{
	addr,err:=net.ResolveUDPAddr("udp","0.0.0.0:53")
	if err != nil {
		log.Println("Failed to resolve the UDP address:",err.Error())
		return err
	}
	conn,err:=net.ListenUDP("udp",addr)
	defer conn.Close()
	if err != nil{
		log.Println("Failed to \"listen to \" the port 53:",err.Error())
		return err
	}
	query := make([]byte,512)
	_,err=conn.Read(query)
	if err != nil{
		log.Println("Failed to read the query from the client:",err.Error())
		_,_=conn.Write([]byte("Failed to read the query .The connection to your device will be closed."))
		conn.Close()
		return err
	}
	reader := bytes.NewReader(conn)
	defer reader.Reset(reader)
	var sendId uint16
	binary.Read(reader,binary.BigEndian,&)
	send()//将结果转发
}

func dnsserverConn(query []byte)error{//打算返回解析值
	/*serverSet:=make([]*net.UDPAddr,5)
	serverSet[0],_=net.ResolveUDPAddr("udp","8.8.8.8:53")
	serverSet[1],_=net.ResolveUDPAddr("udp","1.1.1.1:53")
	serverSet[2],_=net.ResolveUDPAddr("udp","119.29.29.29:53")
	serverSet[3],_=net.ResolveUDPAddr("udp","119.28.28.28:53")
	serverSet[4],_=net.ResolveUDPAddr("udp","1.2.4.8:53")*/
	for i:= 0;i<4;i++{
		conn,err:=net.Dial("udp","8.8.8.8:53")
		if err != nil {
			log.Println("Failed to connect teh DNS server:",err.Error())
			continue
		}
		_,err=conn.Write(query)
		if err != nil {
			log.Println("Failed to send the query to the authority DNS server:",err.Error())
			continue
		}
		answerbuf := make([]byte,512)
		answerLength,err:=conn.Read(answerbuf)
	}
}