package tran_c

//import "CyberLighthouse/dns/process"

type client interface {
	//process.DNSProcess //引入dns
	CreateConn(addr string)
	defaultSend(conn interface{}, message string)
	SendFile(conn interface{}, filename string)
}

type Tcp struct {
}
type Udp struct {
}
