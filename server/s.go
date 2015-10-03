package server

import (
	"config"
	"fmt"
	"io"
	"net"
	"regexp"
)

var (
	addr       string
	port       string
	onlyServer server
	HTTP_r     *regexp.Regexp
)

type server struct {
	flag     bool
	Err      error
	Data     chan [config.MaxNumberOfData][]byte
	listener net.Listener
}

func (s server) GetDataChannel() chan [config.MaxNumberOfData][]byte {
	return s.Data
}

type connHandler interface {
	dealConn(c net.Conn)
}

func NewServer() *server {
	return &onlyServer
}
func init() {
	HTTP_r, _ = regexp.Compile(`^HTTP`)
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
	s.listen()
	if s.Err != nil {
		return s.Err
	}
	go s.accept()
	return nil
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
func (s server) accept() {
	accept(s.listener, s)
}
func (s server) dealConn(c net.Conn) {
	data, err := read(c)
	switch parse(data) {
	case "HTTP":
		dealHTTP(c)
	}
}
func parse(data []byte) string {
	fmt.Println(string(data))
	if isHTTP(data) {
		return "HTTP"
	}
	return "unkown"
}
func isHTTP(data []byte) bool {
	if HTTP_r.Find(data) != nil {
		return true
	}
	return false
}
func read(c net.Conn) ([]byte, error) {
	var result []byte
	buffer := make([]byte, 1024)
	for {
		n, err := c.Read(buffer)
		if err != nil {
			if err != io.EOF {
				dealError(err)
				return nil, err
			}
		}
		result = append(result, buffer[:n]...)
		if n < 1024 {
			break
		}
	}
	return result, nil
}
