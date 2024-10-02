package main

import (
	"fmt"
)

func main() {
	//低配版,输入网址

	fmt.Println("Please enter a DNS server and the port(usually 53):")
	var dnsServer string
	fmt.Scanf("%s",&dnsServer)

	fmt.Println("Please enter which domain address you want to analyse:")
	var domain string
	fmt.Scanf("%s",&domain)

	query,requestLength,duration := Send(dnsServer,domain)//DNS监听端口号常为53
	fmt.Println(query,"\n",requestLength,"\t",duration)
}