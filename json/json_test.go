package json

import (
	"fmt"
	"testing"
)

var test_j = []byte("{'name': 'zhang','age':13 ,'gender' : 'male'}")
var test_j0 = []byte("{\"name\":\"zhang\",\"age\":12,\"gender\":\"female\"}")

func Test_main(t *testing.T) {
	fmt.Println(NewJson(test_j).OutputMap())
	fmt.Println(NewJson(test_j0).OutputMap())
}
