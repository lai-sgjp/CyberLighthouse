package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"CyberLighthouse/tran_c"
)

func main() {

	if len(os.Args) < 2 {
		tran_c.CreateConn() //未指定就默认模式
		return
	}

	//e.g.go run . main -tc="tcp" -a="127.0.0.1.8080" string -c="Hello world!"
	main := flag.NewFlagSet("main", flag.ExitOnError)
	tranptr := main.String("tc", "tcp", "which protocol do you want to use")
	addrptr := main.String("a", "127.0.0.1:8080", "which address do you want to send message")

	if err := main.Parse(os.Args[1:]); err != nil {
		log.Fatal("command err:", err)
	}

	args := main.Args()
	switch args[0] {
	case "string":
		contextptr := flag.NewFlagSet("string", flag.ExitOnError)
		context := contextptr.String("c", "", "what do you want to send")

		if err := contextptr.Parse(args[1:]); err != nil {
			log.Fatal("command err:", err)
		}

		tran_c.Choose(*tranptr, *addrptr, *context)

	case "int":
		contextptr := flag.NewFlagSet("int", flag.ExitOnError)
		context := contextptr.Int("c", 0, "Which number do you want to send?")
		if err := contextptr.Parse(args[1:]); err != nil { //母命令可以有？
			log.Fatal("Failed to read the number you want to send:", err)
		}
		contextEdit := fmt.Sprintf("%d", *context)
		tran_c.Choose(*tranptr, *addrptr, contextEdit)

	case "file":
		fileflag := flag.NewFlagSet("file", flag.ExitOnError)
		filenameptr := fileflag.String("n", "", "which file do you want to send")
		if err := fileflag.Parse(args[1:]); err != nil {
			log.Fatal("Failed to read the file's name:", err)
		}

		filename := fmt.Sprintf("%s", *filenameptr)

		var conn net.Conn
		switch *tranptr {
		case "tcp":
			conn = tran_c.CreateTCPConn(*addrptr)
		case "udp":
			conn = tran_c.CreateUDPConn(*addrptr)
		}

		tran_c.SendFile(conn, filename)
	}

}
