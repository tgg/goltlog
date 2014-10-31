package goltlog

import (
	"errors"
	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
	"github.com/tgg/goltlog/rpc_thrift"
	"fmt"
	"net"
)

func NewLocalLog() (l Log, err error) {
	m := new(logger)
	if err = m.init(); err == nil {
		l = m
	}
	return
}

func NewThriftLogClient(h string, p int) (l Log, err error) {
	var trans thrift.TTransport
	portStr := fmt.Sprint(p)
	if trans, err = thrift.NewTSocket(net.JoinHostPort(h, portStr)); err != nil {
		return
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := rpc_thrift.NewLogClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		return
	}
	l = &thrift_cli{cli: client}
	return
}

func NewThriftLogServer(h string, p int) (l rpc_thrift.Log, err error) {
	var trans thrift.TServerTransport
	portStr := fmt.Sprint(p)
	if trans, err = thrift.NewTServerSocket(net.JoinHostPort(h, portStr)); err != nil {
		
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	handler, err := NewLocalLog()
	if err != nil {
		return
	}

	log := &thrift_srv{h: handler}
	processor := rpc_thrift.NewLogProcessor(l)
	transportFactory := thrift.NewTTransportFactory()
	server := thrift.NewTSimpleServer4(processor, trans, transportFactory, protocolFactory)
	log.srv = server
	l = log
	go server.Serve()
	return
}

func NewLog(method string) (Log, error) {
	switch (method) {
	case "local":
		return NewLocalLog()

	case "thrift":
		if _, err := NewThriftLogServer("localhost", 12345); err == nil {
			return NewThriftLogClient("localhost", 12345)			
		} else {
			return nil, err
		}

	default:
		return nil, errors.New("Unsupported method")	
	}
}
