package rpc

import (
	"fmt"
	"testing"
)

func Test_RPC(t *testing.T) {
	url := "127.0.0.1:55027"
	rpc := NewSeeleRPC(url)
	var result interface{}
	rpc.call("seele.GetInfo", nil, &result)
	fmt.Printf("result is %#v\n", result)
}
