package goltlog

import (
	"errors"
	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
	"github.com/tgg/goltlog/rpc_thrift"
	"fmt"
	"net"
)

func NewThriftLog(h string, p int) (l Log, err error) {
	var trans thrift.TTransport
	portStr := fmt.Sprint(p)
	if trans, err = thrift.NewTSocket(net.JoinHostPort(h, portStr)); err != nil {
		return
	}
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	client := rpc_thrift.NewLogClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		return
	}
	l = &thrift_cli{client: client}
	return
}

func NewLog(method string) (Log, error) {
	switch (method) {
	case "local":
		l := new(logger)
		if err := l.init(); err != nil {
			return nil, err
		}

		return l, nil

	case "thrift":
		l, err := NewThriftLog("localhost", 12345)
		return l, err

	default:
		return nil, errors.New("Unsupported method")	
	}
}
