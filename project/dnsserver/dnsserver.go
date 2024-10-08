package main

import (
	"CyberLighthouse/dns/process"
	"encoding/binary"
	"log"
	"net"
	"os"
)

type trialerror struct {
	message string
}

func (e trialerror) Error() string {
	return e.message
}

type Udp struct {
}

var history []process.Returninfomation

func clientConn() (uint16, []byte, int, net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:53")
	if err != nil {
		log.Println("Failed to resolve the UDP address:", err.Error())
		return 0, nil, 0, net.UDPConn{}, err
	}
	conn, err := net.ListenUDP("udp", addr)
	defer conn.Close()
	if err != nil {
		log.Println("Failed to \"listen to \" the port 53:", err.Error())
		return 0, nil, 0, net.UDPConn{}, err
	}
	query := make([]byte, 512)
	_, err = conn.Read(query)
	if err != nil {
		log.Println("Failed to read the query from the client:", err.Error())
		_, _ = conn.Write([]byte("Failed to read the query .The connection to your device will be closed."))
		conn.Close()
		return 0, nil, 0, *conn, err
	}
	//reader := bytes.NewReader(conn)
	//defer reader.Reset(reader)
	var sendId uint16
	sendId = binary.BigEndian.Uint16(query[:2])
	queryLength := len(query)
	//将结果转发
	return sendId, query, queryLength, *conn, nil
}

func dnsserverConn(query []byte) ([]byte, error) { //打算返回解析值
	/*serverSet:=make([]*net.UDPAddr,5)
	serverSet[0],_=net.ResolveUDPAddr("udp","8.8.8.8:53")
	serverSet[1],_=net.ResolveUDPAddr("udp","1.1.1.1:53")
	serverSet[2],_=net.ResolveUDPAddr("udp","119.29.29.29:53")
	serverSet[3],_=net.ResolveUDPAddr("udp","119.28.28.28:53")
	serverSet[4],_=net.ResolveUDPAddr("udp","1.2.4.8:53")*/
	for i := 0; i < 4; i++ {
		conn, err := net.Dial("udp", "8.8.8.8:53")
		if err != nil {
			log.Println("Failed to connect teh DNS server:", err.Error())
			continue
		}
		_, err = conn.Write(query)
		if err != nil {
			log.Println("Failed to send the query to the authority DNS server:", err.Error())
			continue
		}
		answerbuf := make([]byte, 512)
		_, err = conn.Read(answerbuf)
		if err != nil {
			log.Println("Failed to receive the response from the DNS server.")
		}
		return answerbuf, nil
	}
	log.Println("Failure.")
	err := trialerror{message: "Have tried 5 times but each of it has failed."}
	return nil, err
}

func storge(answer process.Returninfomation) error {
	err := append(history, answer)
	if err != nil {
		log.Println("Failed to storge the query and the record this time.")
		return answer.Err
	}
	return nil
}

func main() {
	sendId, query, queryLength, conn, err := clientConn()
	if err != nil {
		os.Exit(1)
	}
	answerbuf, err := dnsserverConn(query)
	result := <-process.Udp.Parse(queryLength, sendId, conn)

}
