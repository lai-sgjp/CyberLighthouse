package tran

type Tcp interface {
	CreateTCPSer(ports []string)
	ProcessTCP(conn net.Conn)
	textMode(conn net.Conn, buf string) error
	fileMode(conn net.Conn)
}

type Udp interface {
	CreateUDPSer(ports []string)
	ProcessUDP(conn *net.UDPConn)
	textModeudp(conn *net.UDPConn, newBuf string) error
	fileModeudp(conn *net.UDPConn)
}
