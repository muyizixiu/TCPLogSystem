package server

import (
	"config"
	"fmt"
	"io"
	"net"
	"regexp"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	addr       string
	port       string
	onlyServer server
	HTTP_r     *regexp.Regexp
)

type server struct {
	flag     bool
	Err      error
	Data     chan []byte
	listener net.Listener
}

func (s server) GetDataChannel() chan []byte {
	return s.Data
}

type connHandler interface {
	dealConn(c net.Conn)
}

func NewServer() *server {
	return &onlyServer
}
func init() {
	initData()
	HTTP_r, _ = regexp.Compile(`HTTP[^(\r\n)]+\r\n`)
	onlyServer = server{flag: true, Data: make(chan []byte, config.MaxNumberOfData)}
}
func initData() {
	addr = ""
	port = "80"
}
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

type HTTPHandler interface {
	dealHTTP()
}

func (s server) dealConn(c net.Conn) {
	con := &Conn{conn: c}
	data, err := con.read()
	if err != nil {
		s.Err = err
		con.Close()
		return
	}
	switch parse(data) {
	case "HTTP":
		h := HTTPHandler(newHTTP(con))
		h.dealHTTP()
	}
	c.Close()
}
func parse(data []byte) string {
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

type Conn struct {
	conn   net.Conn
	buffer []byte
}

func (c *Conn) readAll() error {
	buf, err := read(c.conn)
	if err != nil {
		return err
	}
	c.buffer = append(c.buffer, buf...)
	return nil
}
func (c *Conn) read() ([]byte, error) {
	buf, err := read(c.conn)
	if err != nil {
		return nil, err
	}
	c.buffer = append(c.buffer, buf...)
	return buf, nil
}
func (c Conn) Close() {
	c.conn.Close()
}
