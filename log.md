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
    >一个文件夹里面只能有一个包
    >mod是一个module，里面有很多package，相互引用路径从module开始

---

## 2024.10.02
- 0200:还没有解决关于命令行参数的bug，但我想先摆烂睡觉zzz
- 0911:睡醒了，开始解决
- 0935:解决了命令行参数的报错（主要是关于命令行参数该如何使用，如何选取特定的命令行参数，
- 对代码进行调试（完善错误处理不要无限输出等）
- 1646跳过任务1的get,post和限制访问频率的子项，进入任务二
    - DNS协议的作用
        >通过域名查出对应的IP地址，从而建立通信
    - DNS查询流程
        >依次查询
        >在终端里面输入`dig <网址>`即可查到，分为六个部分
        >1.查询的参数和各自统计
        >2.查询的内容`QUESTION SECTION`
        >3.DNS服务器的答复`ANSWER SECTION`
        >4.NS记录，即哪些服务器管理该网址的DNS记录`AUTHORITY SECTION`
        >5.在上一项中的服务器IP地址`ADDITIONAL SECTION`
        >6.这次查询的信息
        >不过在我电脑上只看到1，2，3，6（尴尬）
    - DNS服务器的分类与作用
        >内外网分
        >内网DNS服务器
        >1.解析**内网**的私有IP地址
        >2.缓存位于内网的计算机访问外网的DNS查询结果
        >3.提供负载均衡作用，将请求分配给不同的处理器
        >4.安全防控（异想天开：GFW原理之一会不会就是此？）
        >
        >外网DNS服务器
        >1.全球访问，齐全无遗漏。
        >
        >
        >或者
        >根域名服务器:查询顶级域名`TLD`的DNS
        >TLD服务器:查询各个一级域名的DNS
        >权威服务器:就是一级域名的服务器，往下一次递推
        >递归域名服务器:将上述自动化
        >
        >可以指定DNS服务器:`dig @<DNS服务器地址> <网址>`
    - DNS报文格式，各标志位的含义
        1.大体结构:**头部-问题-答案-授权信息-附加信息**
        2.头部:标识（ID），标志，问题、答案、授权信息、附加信息等的数量（count）
            >标志中有：
            >QR（0表示查询请求，即是需求端；1表示查询应答，即是服务端）
            >操作码（0表示标准查询；1表示反向查询即根据IP地址找域名(据说可以看看邮件发送的ip是否属实)；2表示服务器状态请求）
            >AA权威回答；TC截断（512bytes）；RD期望递归（递归域名服务器，1为TRUE）；RA可递归，告知client可以进行递归查询；保留字段（为了可能的扩展）
            >响应码（0为无差错，3为有差错from权威服务器）
        3.问题:待查询域名，查询类型（DNS记录类型），类(通常为1，表示TCP/IP地址)
            >域名：每级域名前会加上一位表示该级的长度
        4.答案:被查询域名，查询类型，类（前三个与问题一致），有效期(TTL，即在多长时间内无需进行DNS查询)，数据长度，数据
        5.授权信息:提示哪些服务器可以有权限去查询，比如说各个国籍
        6.其他注意:域名压缩(第一个字节最高两位都是1，剩余与第二个字节组成偏移量以此推断出第一次出现时哪些部分是域名)，节省资源
        
        >其他
        >DNS记录类型:`A``NS``MX``CNAME``PTR`，分别是*地址*，*域名服务器记录*，*邮件记录*，*跳转记录*，*逆向查询记录*
        >`dig <类型> <网址>`可以指定记录类型
        >`dig -x <ip>`可以反向查询
    - 分级查询:每一级域名都有自己的NS记录,从顶级域名到后面一个一个的指路
        >`dig +trace <网址>`显示分级查询过程(不过貌似实操时被GFW封住了？显示超时)
- 2050学习如何用`go`实现dns的解析
    >参考:*https://www.cnblogs.com/chase-wind/p/6814053.html?utm_campaign=studygolang.com&utm_medium=studygolang.com&utm_source=studygolang.com*
- 0130:成功实现DNS发送请求报文，但还差解析响应报文
    >问题：向8.8.8.8:53发送B站和Bing的两次一样，且后面全为0.而用1.1.1.1:53则是前一个问题得到解决
## 2024.10.03
- 完成task3!
## 学习路线
1. 不会的先想一想->bing->转至各大网站
2. `go`相关：*go bible*,*go语言中文网*等
3. 网络相关：阮一峰的网络日志，小菜学网络

## 如何测试
1. 按照功能一个个输入
2. go语言有问题就按照错误提示改
3. 没有任何报错达不到预期打印日志

## 心路历程
1. 学习时：有趣，充满干劲
2. 调错时：en看懂了错误->nani？竟然不行->暴躁->听歌缓解->继续调试->成功兴奋
3. 临近ddl：完了完了怎么还没做完已经要上交了？？？nani？8号要交英语presentation还没准备？？？ -> 听歌缓解 -> 继续


## 实现的功能
1. task1中的*基础任务*（还没有时间写发送整数类型，但是有思路：`binary.write(conn,message,)），*进阶任务、支持文件(有bug,还未调解成功：位置在于client端发送的文件大小server端接收不到)、并发处理*
2. task2
3. task3中的*基础任务*，*进阶任务（仅支持AAAA，NS等由于还未搞清楚域名压缩所以未成功）*
4. task4中的*基础任务（命令行参数部分实现）*，*进阶任务的AAAA*
5. task5与NS等任务无关的有思路但没时间（内存缓存开切片，循环发送用循环+计时器通道，硬盘缓存最初想法是文件存储，看到数据库时已经要提交了没考虑）

## 补充日志（git提交记录）
Date:   Sun Oct 6 08:58:03 2024 +0800

    set up the DNS imtermidiate server primarily.

commit 
Author: Lai 
Date:   Sun Oct 6 02:05:53 2024 +0800

    add the function to remind successful result and try to send the file which is a little bugs--Don't send the size of the file

commit 
Author: Lai 
Date:   Sat Oct 5 20:50:48 2024 +0800

    correct a mistake of 'net.Conn'vs'*net.UDPConn' and add the function to provide DNS service with the TCP protocol.

commit (dev)
Author: Lai 
Date:   Sat Oct 5 19:58:09 2024 +0800

    Correct some bugs in the TCP and UDP establishment and the data reception.

commit 
Author: Lai 
Date:   Sat Oct 5 13:08:28 2024 +0800

    remove some test files

commit 
Author: Lai 
Date:   Sat Oct 5 13:06:56 2024 +0800

    combine dns service and the client

commit 
Author: Lai 
Date:   Sat Oct 5 11:57:02 2024 +0800

    use the interface

commit 
Author: Lai 
Date:   Sat Oct 5 01:06:59 2024 +0800

    Correct the format of the DNS response

commit 
Author: Lai 
Date:   Fri Oct 4 19:35:05 2024 +0800

    Change the way to store the DNS report

commit
Author: Lai 
Date:   Fri Oct 4 14:07:06 2024 +0800

    correct the usage of flags in client.go and the way to read buffer in the tcp protocol

commit 
Author: Lai 
Date:   Thu Oct 3 22:22:58 2024 +0800

    remove some comments and update the log.md

commit 
Author: Lai <
Date:   Thu Oct 3 22:14:27 2024 +0800

    debug

commit 
Author: Lai 
Date:   Thu Oct 3 20:47:57 2024 +0800

    generate random id in query

commit 
Author: Lai 
Date:   Thu Oct 3 20:44:17 2024 +0800

    debug the query part and change the encoding way from LittleEndian to BigEndian

commit 
Author: Lai 
Date:   Thu Oct 3 17:05:58 2024 +0800

    add the function to analyse the response code and reorganise the module

commit 
Author: Lai 
Date:   Thu Oct 3 12:11:02 2024 +0800

    add the function of parse dns(did not test)

commit 
Author: Lai 
Date:   Thu Oct 3 02:00:07 2024 +0800

    commit the newest log

commit 
Author: Lai 
Date:   Thu Oct 3 01:54:10 2024 +0800

    add the function to send query for the dns

commit 
Author: Lai 
Date:   Wed Oct 2 20:30:19 2024 +0800

    finish task 2

commit 
Author: Lai 
Date:   Wed Oct 2 16:43:33 2024 +0800

    debug for the 2nd time

commit 
Author: Lai 
Date:   Wed Oct 2 16:14:32 2024 +0800

    use modules to organise the code for the 2nd time

commit 
Author: Lai 
Oct 2 15:58:38 2024 +0800

    debug for the first time

commit 
Author: Lai <
Date:   Wed Oct 2 10:31:05 2024 +0800

    add the function to send files

commit 
Author: Lai 
Date:   Wed Oct 2 09:33:52 2024 +0800

    use the parameter in the command line

commit 
Author: Lai 
Date:   Tue Oct 1 21:50:36 2024 +0800

    use modules to organize the codes

commit 
Author: Lai 
Date:   Tue Oct 1 18:59:57 2024 +0800

    wrote a client.go program

commit 
Author: Lai 
Date:   Tue Oct 1 16:17:21 2024 +0800

    In this record I successfully run the first server.go program,and I commit my learning record on the Oct.1st as well.