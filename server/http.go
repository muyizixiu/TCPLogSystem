package server

import (
	"errors"
	"net"
	"regexp"
)

var httpHead_r *regexp.Regexp

func init() {
	httpHead_r, _ = regexp.Compile(`^HTTP[^(\r\n\r\n]+\r\n\r\n`)
}

type HTTP struct {
	ip     string
	head   []byte
	origin []byte
	data   json
	conn   net.Conn
	err    error
}

func newHTTP(c net.Conn) *HTTP {
	return &HTTP{ip: c.RemoteAddr().String(), conn: c}
}
func (h *HTTP) dealHTTP() {
	go h.readData()
	h.reply()
	h.readHead()
	if h.err != nil {
		return
	}
	(newJson(h.origin[len(h.head):], h.ip)).send()
}
func (h *HTTP) close() {
	h.conn.Close()
}
func (h *HTTP) readData() {
	var buffer []byte
	part, err := read(h.conn)
	if err != nil {
		h.close()
		h.err = err
		return
	}
	h.origin = append(buffer, part...)
}
func (h *HTTP) reply() {
	h.conn.Write([]byte("HTTP/1.1 200 OK\r\nData: Sun, 04 Oct 2015 09:07:00\r\n\r\n"))
	h.close()
}
func (h *HTTP) readHead() {
	h.head = httpHead_r.Find(h.origin)
	if h.head == nil {
		h.err = errors.New("fake http format")
	}
}
