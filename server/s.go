package server

import (
	"fmt"
	"net"
)

var (
	addr       string
	port       string
	onlyserver server
)

type server struct {
	flag     bool
	Err      error
	Data     chan [config.MaxNumberOfData][]byte
	listener net.Listener
}
type connHandler interface {
	dealConn(c net.Conn)
}

func NewServer() {
	return &onlyServer
}
func init() {
	onlyServer = server{flag: true, Data: make(chan [config.MaxNumberOfData][]byte)}
}
func initData() {}
func checkErr(err error) bool {
	if err != nil {
		dealError(err)
		return true
	}
	return false
}
func dealError(err error) {
	fmt.Println(err)
}
func (s server) start() error {
}
func listen() (net.Listener, error) {
	l, err := net.Listen("tcp", addr+":"+port)
	if checkErr(err) {
		return nil, err
	}
	return l, nil
}
func (s *server) listen() {
	s.listener, s.Err = listen()
}
func accept(l net.Listener, h connHandler) {
	for {
		c, err := l.Accept()
		if checkErr(err) {
			continue
		}
		go h.dealConn(c)
	}
}
