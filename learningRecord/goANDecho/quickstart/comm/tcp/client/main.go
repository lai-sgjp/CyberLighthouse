package main
import (
	"http/net"
	"fmt"
	"os"
	"bufio"
)

func main() {
	conn,err := net.Dial("tcp","127.0.0.1:8080")//客户端建立连接，前面的是传输协议，后面的是监听的网址和端口
	if err != nil {
		fmt.Println("err:" err)
		return 
	}//每次建立连接，监听器等都要进行返回值err的检查

	defer conn.Close()//提前设置好连接关闭
	inputReader := bufio.NewReader(os.Stdin)//建立读取缓存的阅读器(类似scanf?)
	for {
		input ,_ := inputReader.ReadString("\n")//读取用户输入，读到换行为止（包括换行符）
		inputInfo := strings.Trim(input,"\n")//对用户输入进行格式化，去除所有的`\n``\r`
		if strings.ToUpper(inputInfo) == "Q" {//将所有字母转为大写，如果输入了`q`或者`Q`就会退出程序
			return
		}
	}

	_,err := conn.Write([]byte(inputInfo))//发送数据
	if err!= nil {
		return
	}

	buf := [512]byte{}
	n,err := conn.Read(buf[:])

	if err != nil {
		fmt.Println("recv failed,err:",err)
		return
	}

	fmt.Println(string(buf[:n]))
}

//建立连接`net.Dial`->读取数据->格式化数据->发送数据`conn.Write`->读取数据`conn.Read`->打印数据 ->关闭连接