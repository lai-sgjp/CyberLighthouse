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
	if input == "" {
		log.Println("You enter nothing.We will use the port \"8080\" by default.")
		input = "8080"
	}
	ports := strings.Split(input, " ")
	return ports, nil
}

func main() {
	ports, err := getport()
	if err != nil {
		log.Fatal("Failed to get the port:", err)
	}

	fmt.Println("Whether turn on the UDP service or not(y/n)")
	var choice string
	fmt.Scanln(&choice)
	choice = strings.ToLower(strings.TrimSpace(choice))
	if choice == "y" {
		u := tran.Udp{} //创建实例(不在一个包里面但又要使用接口里面的方法)
		go u.CreateSer(ports)
	} else if choice != "y" && choice != "n" {
		log.Println("Input error!We will only use the TCP service by default.")
	}
	t := tran.Tcp{}
	go t.CreateSer(ports)

	select {} //保持server端长时间不关闭
}
