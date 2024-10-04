package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"CyberLighthouse/tran"
)

func getport() ([]string, error) {
	fmt.Println("Please enter port seperate by space.(e.g.\"8000 8001\" \"8002\")")
	//获得一连串的输入
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n')
	if err != nil {
		return nil, err
	}
	//格式化输入
	input = strings.TrimSpace(input)
	ports := strings.Split(input, " ")
	return ports, nil
}

func main() {
	ports, err := getport()
	if err != nil {
		log.Fatal("Failed to get the port:", err)
	}

	fmt.Println("Whether turn on the UDP service or not(y/n)")
	var choice rune
	fmt.Scanln("%v", &choice)
	if choice == 'y' {
		go tran.CreateUDPSer(ports)
	} else {
		go tran.CreateTCPSer(ports)
	}

	select {} //保持server端长时间不关闭
}
