package main

import (
	"CyberLighthouse/dns/report"
	"fmt"
	"log"
	"time"
)

func main() {
	fmt.Println("Please enter a DNS server and the port(port usually is 53):") //期待加上超时
	var dnsServer string
	fmt.Scanf("%s", &dnsServer)

	fmt.Println("Please enter which domain address you want to analyse:")
	var domain string
	fmt.Scanf("%s", &domain)

	query, sendId, queryLength, conn, duration, err := report.Send(dnsServer, domain)
	if err != nil {
		log.Fatal("Failed to send the query to  the DNS server:", err.Error())
	}
	fmt.Println(query, "\n", queryLength, "\t", duration)

	//设置计时器
	timeout := 5 * time.Second
	timeoutChan := time.After(timeout)
	select {
	case <-timeoutChan:
		fmt.Println("Wait for the response for too long...")
		conn.Close()
		return
	case returninfo := <-report.Parse(queryLength, sendId, conn):
		if returninfo.Err != nil {
			log.Fatal("Parse Error:", returninfo.Err.Error())
		}
		//注意：大写才可以被导出至不同的包！！
		//困扰很久的报错：returninfo.other undefined (type report.returninfomation has no field or method other)
		fmt.Printf("{%v}\n", returninfo.Other.Bytes())
		conn.Close()
	}
}
