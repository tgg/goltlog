// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package rpc_thrift

import (
	"fmt"
	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

type Log interface {
	LogProducer
}

type LogClient struct {
	*LogProducerClient
}

func NewLogClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *LogClient {
	return &LogClient{LogProducerClient: NewLogProducerClientFactory(t, f)}
}

func NewLogClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *LogClient {
	return &LogClient{LogProducerClient: NewLogProducerClientProtocol(t, iprot, oprot)}
}

type LogProcessor struct {
	*LogProducerProcessor
}

func NewLogProcessor(handler Log) *LogProcessor {
	self165 := &LogProcessor{NewLogProducerProcessor(handler)}
	return self165
}

// HELPER FUNCTIONS AND STRUCTURES