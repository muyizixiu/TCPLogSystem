package server

import (
	"errors"
	"sync"
)

const (
	bufferSize = 2048
)

var (
	lengthErr    = errors.New("wrong length")
	checkCode    = byte(12)
	checkCodeErr = errors.New("content error")
	MaxLengthErr = errors.New("out of MaxLength")
)

type TCP struct {
	ip         string
	head       []byte
	buffer     []byte
	bufferLock *sync.Mutex
	conn       *Conn
	err        error
}

func newTCP(c *Conn) *TCP {
	return &TCP{ip: c.conn.RemoteAddr().String(), conn: c, bufferLock: &sync.Mutex{}}
}
func (t *TCP) dealTCP() {
	for {
		t.readData()
		j := t.getData()
		if t.err != nil {
			if t.err == lengthErr {
				continue
			}
			t.dealErr()
			return
		}
		j.send()
	}
}
func (t *TCP) checkHead() {
}
func (t *TCP) responseHead() {
}
func (t *TCP) getData() *json {
	if len(t.buffer) < 2 {
		t.err = lengthErr
		return nil
	}
	length := int(t.buffer[0])*256 + int(t.buffer[1])
	if len(t.buffer) < (length + 2) {
		t.err = lengthErr
		return nil
	}
	if t.buffer[length+1] != checkCode {
		t.err = checkCodeErr
		return nil
	}
	re := t.buffer[2 : length+2]
	t.bufferLock.Lock()
	t.buffer = t.buffer[length+2:]
	t.bufferLock.Lock()
	return newJson(re, t.ip)
}
func (t *TCP) readData() {
	data, err := t.conn.read()
	if err != nil {
		t.err = err
		t.close()
	}
	if len(data) > bufferSize {
		t.err = MaxLengthErr
		return
	}
	t.bufferLock.Lock()
	t.buffer = append(t.buffer, data...)
	if len(t.buffer) > bufferSize {
		t.err = MaxLengthErr
	}
	t.bufferLock.Unlock()
}
func (t *TCP) close() {
	t.conn.Close()
}
func (t *TCP) checkErr() bool {
	if t.err != nil {
		t.dealErr()
		return true
	}
	return false
}
func (t *TCP) dealErr() {
	t.conn.Close()
}
