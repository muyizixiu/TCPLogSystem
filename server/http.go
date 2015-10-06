package server

import (
	"errors"
	"regexp"
)

var httpHead_r *regexp.Regexp

func init() {
	httpHead_r, _ = regexp.Compile(`[(POST)(OPTIONS)] / HTTP/1.1[\s\S]*?(\r\n){2,}`)
}

type HTTP struct {
	ip     string
	head   []byte
	origin []byte
	data   json
	conn   *Conn
	err    error
}

func newHTTP(c *Conn) *HTTP {
	return &HTTP{ip: c.conn.RemoteAddr().String(), conn: c}
}
func (h *HTTP) dealTCP() {
	h.dealHTTP()
}
func (h *HTTP) dealHTTP() {
	h.origin = h.conn.buffer
	h.reply()
	h.readHead()
	if h.err != nil {
		dealError(h.err)
		return
	}
	(newJson(h.origin[len(h.head):], h.ip)).send()
}
func (h *HTTP) close() {
	h.conn.Close()
}
func (h *HTTP) readData() {
	var buffer []byte
	err := h.conn.readAll()
	part := h.conn.buffer
	if err != nil {
		h.close()
		h.err = err
		return
	}
	h.origin = append(buffer, part...)
}
func (h *HTTP) reply() {
	h.conn.conn.Write([]byte("HTTP/1.1 200 OK\r\nData: Sun, 04 Oct 2015 09:07:00\r\n\r\n"))
	h.close()
}
func (h *HTTP) readHead() {
	h.head = httpHead_r.Find(h.origin)
	if h.head == nil {
		h.err = errors.New("fake http format")
	}
}
