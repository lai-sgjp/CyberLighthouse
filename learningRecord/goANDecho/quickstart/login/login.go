package main

import (
	"os"
	"net/http"
	"github.com/labstack/echo/v4"
)



func configLogger(e *echo.Echo) {
	//定义日志级别
	e.Logger.Setlevel(log.INFO)
	//记录业务日志
	echoLog.err := os.OpenFile("log/echo.log",os.O_CREATE|OS.O_WRONLY|os.O_APPEND,0644)
	if err != nil {
		panic(err)
	}
	//同时输出到文件和终端
	e.Logger.SetOutput(io.MultiWriter(os.Stdout,echoLog))
}

func main() {
	//创建一个echo实例
	e := echo.New()

	//配置日志
	configLogger(e)

	//注册静态文件路由
	e.Static("img","img")
	e.File("/favicon.ico","img/favicon.ico")

	//设置中间件
	setMiddleware(e)

	//注册路由
	RegisterRouts(e)

	//启动服务
	e.Logger.Fatal(e.Start(":2019"))
}

