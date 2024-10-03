package main

import (
	"CyberLighthouse/dns/report"
	"fmt"
	"log"
)

func main() {
	//低配版,输入网址

	fmt.Println("Please enter a DNS server and the port(usually 53):")
	var dnsServer string
	fmt.Scanf("%s", &dnsServer)

	fmt.Println("Please enter which domain address you want to analyse:")
	var domain string
	fmt.Scanf("%s", &domain)

	query, sendId, queryLength, conn, duration, err := report.Send(dnsServer, domain) //DNS监听端口号常为53
	fmt.Println(query, "\n", queryLength, "\t", duration)
	if err != nil {
		log.Fatal("Failed to send the query to  the DNS server:", err.Error())
	}
	otherbuf, err := report.Parse(queryLength, sendId, conn)
	if err != nil {
		log.Fatal("Parse Error:", err.Error())
	}
	fmt.Printf("{%v}\n", otherbuf)
}
