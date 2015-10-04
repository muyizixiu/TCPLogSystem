package server

import ()

type json struct {
	data        []byte
	ip          string
	dataChannel chan []byte
}

func newJson(d []byte, ipv string) *json {
	return &json{data: d, dataChannel: onlyServer.GetDataChannel(), ip: ipv}
}
func (j json) send() {
	j.dataChannel <- j.data
}
