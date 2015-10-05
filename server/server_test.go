package server

import (
	"fmt"
	"runtime"
	"testing"
)

func init() {
	runtime.GOMAXPROCS(2)
}
func Test_main(t *testing.T) {
	onlyServer.start()
	data := onlyServer.GetDataChannel()
	for {
		fmt.Println(string(<-data))
	}
}
