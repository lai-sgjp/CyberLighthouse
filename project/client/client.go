package main

import (
	"flag"
	"os"
	"log"
	"fmt"
	"net"

	"example2.com/createConn"
)

func main() {

	if len(os.Args) < 2 {
		tran_c.CreateConn()//未指定就默认模式
		return
	}

	main := flag.NewFlagSet("main",flag.ExitOnError)
	tranptr := main.String("tc","tcp","which protocol do you want to use")//tranconfig//这里是直接输入到命令行，无需加上tag(无索引为1)，直接用指针进行操作
	addrptr := main.String("a","127.0.0.1:8080","which address do you want to send message")//无index:2，直接用指针进行操作
	//if err := flag.Parse(); err != nil {
	//	log.Fatal(err)
	//	//1是代表整个命令参数（包括参数）
	//}//flag.Parse()是没有返回值的
	if err := main.Parse(os.Args[1:]); err != nil {
			log.Fatal("command err:",err)
		}

	args := main.Args()//将母命令后面的参数重新视为一个切片
	switch args[0] {//有子行的必须放在第一个
	case "string":
		contextptr := flag.NewFlagSet("string",flag.ExitOnError)
		context := contextptr.String("c","","what do you want to send")
		//if err := flag.Parse();err != nil {
		//	log.Fatal("Failed to read the message you want to send:",err)
		//}

		//context = *context
		if err := contextptr.Parse(args[1:]); err != nil {
			log.Fatal("command err:",err)
		}

		tran_c.Choose(*tranptr,*addrptr,*context)


	case "int":
		contextptr := flag.NewFlagSet("int",flag.ExitOnError)
		context := contextptr.Int("c",0,"Which number do you want to send?")
		if err := contextptr.Parse(args[1:]);err != nil {//母命令可以有？
			log.Fatal("Failed to read the number you want to send:",err)
		}
		contextEdit := fmt.Sprintf("%d",*context)
		tran_c.Choose(*tranptr,*addrptr,contextEdit)


	case "file":
		fileflag := flag.NewFlagSet("file",flag.ExitOnError)
		filenameptr := fileflag.String("n","","which file do you want to send")
		if err := fileflag.Parse(args[1:]);err != nil {
			log.Fatal("Failed to read the file's name:",err)
		} 

		filename := fmt.Sprintf("%s",*filenameptr)

		var conn net.Conn
		switch *tranptr {
		case "tcp":
			conn = tran_c.CreateTCPConn(*addrptr)
		case "udp":
			conn = tran_c.CreateUDPConn(*addrptr)//给已经存在的变量修改值应该直接使用`=`而非`:=`
		}

		tran_c.SendFile(conn,filename)
	}
	
/*
	str := flag.NewFlagSet("string",flag.ExitOnError)
	tranPtr := flag.String("tran","tcp","the protocol")
	strptr := str.String("context","","what do you want to send")
	

	integer := flag.NewFlagSet("int",flag.ExitOnError)
	tranP := flag.String("tran","tcp","the protocol")
	intptr := integer.Int("context",0,"which integer do you want to send")
	//输入的内容应该放在最后面
*/

/*
	switch os.Args[1] {
	case "string":
		str.Parse(os.Args[2:])
		if *tranPtr != "udp" {
			conn := tran_c.CreateTCPConn(os.Args[4])
			_,err := conn.Write([]byte(*strptr))
			if err != nil {
				log.Fatal("Failed to send the message:",err)
			}
			break
		}
		
	case "int":
		integer.Parse(os.Args[2:])
		if *tranP != "udp" {
			conn := tran_c.CreateTCPConn(os.Args[4])
			buf := fmt.Sprintf("%d",*intptr)//简洁将整数转为字符串，因为go中只允许发送utf-8或者二进制编码
			_,err := conn.Write([]byte(buf))
			if err != nil {
				log.Fatal("Failed to send the message:",err)
			}
			break
		}

		
	}
*/
	//tran_c.CreateConn()
}
