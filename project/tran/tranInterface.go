package tran

type Tcp struct {
	/*
		addr  string
		ports []string
	*/
}

type Udp struct {
	/*
		addr  string
		ports []string
	*/
}

type Server interface {
	CreateSer(ports []string)
	Process(conn interface{})
	textMode(conn interface{}, newBuf string)
	fileMode(conn interface{})
}
