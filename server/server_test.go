package server

import (
	"fmt"
	"testing"
)

func Test_main(t *testing.T) {
	onlyServer.start()
	data := onlyServer.GetDataChannel()
	for {
		fmt.Println(<-data)
	}
}
