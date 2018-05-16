/**
*  @file
*  @copyright defined in monitor-api/LICENSE
 */

package rpc

import (
	"net/rpc/jsonrpc"

	"github.com/seeleteam/monitor-api/core/logs"
)

type logger interface {
	Println(v ...interface{})
}

// MonitorRPC json_rpc client
type MonitorRPC struct {
	url    string
	scheme string
	Debug  bool
}

// New create new json_rpc client with given url
func newRPC(url string, options ...func(rpc *MonitorRPC)) *MonitorRPC {
	rpc := &MonitorRPC{
		url:    url,
		scheme: "tcp",
	}
	for _, option := range options {
		option(rpc)
	}
	return rpc
}

func NewSeeleRPC(url string, options ...func(rpc *MonitorRPC)) *MonitorRPC {
	return newRPC(url, options...)
}

func (rpc *MonitorRPC) call(serviceMethod string, args interface{}, reply interface{}) error {
	conn, err := jsonrpc.Dial(rpc.scheme, rpc.url)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		return err
	}

	err = conn.Call(serviceMethod, args, &reply)
	if err != nil {
		return err
	}
	if rpc.Debug {
		logs.Debug("%s\nRequest: %v\nResponse: %v\n", serviceMethod, args, &reply)
	}
	return nil
}
