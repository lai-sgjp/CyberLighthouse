# 赖浩君的国庆实习日志 

## 2024.09.30
**开始征程！**
1.20点44分：学习如何用`go`实现一个简单的Web服务端（学习`go + echo`框架）
    >学习网站：*https://echo.labstack.com/docs/quick-start*

---

## 2024.10.01
- 0046:终于将`echo`框架成功安装在`wsl2`中
    >解决方法:在`wsl2`显示设置`go env -w GOPROXY="https://goproxy.io"``go env -w GO111MODULE="on"`之后再`go get -u github.com/labstack/echo/v4`
    >`wsl2`与`windows`不是一个系统
- 0200:睡觉
- 0802:学习如何使用`golang`实现`TCP`通信和`UDP`通信(通过`net`包)
    >https://zhuanlan.zhihu.com/p/302547547
- 0841:学习`echo`中的**Customization**
    >https://echo.labstack.com/docs/customization
- 0931:打算先用`net`先建立起server,client
    >因为没找到`echo`如何指定传输协议是`tcp`还是`udp`
- 1600:终于让第一个`go`写出来的`server`端运行起来了(QwQ)
- 1854:`client`端也成功了！而且两者可以互相通信！！！（先完成再说）
- 2125:解决`go`模块的引用问题

---

## 2024.10.02
- 0200:还没有解决关于命令行参数的bug，但我想先摆烂睡觉zzz
- 0911:睡醒了，开始解决
- 0935:解决了命令行参数的报错（主要是关于命令行参数该如何使用，如何选取特定的命令行参数，
- 对代码进行调试（完善错误处理不要无限输出，