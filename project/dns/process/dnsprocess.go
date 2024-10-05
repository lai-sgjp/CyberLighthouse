package process

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

/*
//这里没有想好到底哪里需要多态性的接口
type DNSProcess interface {
	DNS(dnsServer, domain string)
}
type Udp struct {
}
*/
//func (u *Udp) DNS(dnsServer, domain string) {
func DNS(dnsServer, domain string) {
	_, sendId, queryLength, conn, duration, err := Send(dnsServer, domain)
	if err != nil {
		log.Fatal("Failed to send the query to  the DNS server:", err.Error())
	}
	//fmt.Println(query, "\n", queryLength, "\t", duration)
	fmt.Println(queryLength, "\t", duration)

	//设置计时器
	timeout := 5 * time.Second
	timeoutChan := time.After(timeout)
	select {
	case <-timeoutChan:
		fmt.Println("Wait for the response for too long...")
		conn.Close()
		return
	case returninfo := <-Parse(queryLength, sendId, conn):
		if returninfo.Err != nil {
			log.Fatal("Parse Error:", returninfo.Err.Error())
		}
		//注意：大写才可以被导出至不同的包！！
		//困扰很久的报错：returninfo.other undefined (type report.returninfomation has no field or method other)
		//fmt.Printf("{%v}\n", returninfo.Other.Bytes())

		//executablePath, err := os.Executable()
		if err != nil {
			log.Fatal("Failed to create the file because failed to get the code location:", err.Error())
		}
		//execDir := filepath.Dir(executablePath) //获取该文件的绝对位置
		domain = strings.Replace(domain, ".", "", -1)
		//strings.Trim表示去除首位两端指定的字符
		//strings.Replace表示去除替换某些字符，-1表示全部应用，前一个字符串表示要被替换的字符，后一个表示替换为什么字符
		//relativePath := "dnsReport/" + domain + ".txt"       //定义相对位置
		//absolutePath := filepath.Join(execDir, relativePath) //组合成绝对位置
		file, err := os.Create(domain + "DNSReport.json")
		if err != nil {
			log.Fatal("File creation failed:", err)
		}
		defer file.Close()
		dnsinfo, _ := json.Marshal(returninfo)
		_, err = file.Write(dnsinfo) //创建写入（位置，内容，权限）
		//不能直接写入字符串，但可以写入json,bin,xml
		if err != nil {
			log.Println("Failed to storge the data in the file:", err.Error())
		}
		fmt.Println("Data is successfully stored in the " + domain + "DNSReport.json")

		conn.Close()
	}
}
