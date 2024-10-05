package main

//必须为main不然不能运行
//因为库里面不能有main()函数不然就有多个执行入口

import (
	"CyberLighthouse/dns/process"
	"fmt"
	"log"
	"os"

	//"path/filepath"
	"strings"
)

type DNSProcess interface {
	DNS(dnsServer, domain string)
}

func main() {
	fmt.Println("Please enter a DNS server and the port(port usually is 53):") //期待加上超时
	var dnsServer string
	fmt.Scanf("%s", &dnsServer)
	if strings.Replace(dnsServer, " ", "", -1) == "" {
		log.Println("Since you enter nothing/break, we will use 8.8.8.8:53 by default.")
		dnsServer = "8.8.8.8:53"
	}

	fmt.Println("Please enter which domain address you want to analyse:")
	var domain string
	fmt.Scanf("%s", &domain)
	if strings.Replace(dnsServer, " ", "", -1) == "" {
		log.Println("Since you enter nothing/break, we will exit the pocess.")
		os.Exit(1)
	}
	//u := process.Udp{}
	//u.DNS(dnsServer, domain)
	fmt.Println("Which way do you want to choose?Please enter \"udp\" or \"tcp\"")
	var choice string
	fmt.Scanf("%s", &choice)
	if choice != "udp" && choice != "tcp" {
		log.Println("You enter protocol we don't support..We will use \"udp\" as default")
		choice = "udp"
	}
	process.DNS(dnsServer, domain, choice)
}
