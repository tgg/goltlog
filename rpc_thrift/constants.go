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

var SECURITY_ALARM U16
var FAILURE_ALARM U16
var DEGRADED_ALARM U16
var EXCEPTION_ERROR U16
var FLOW_CONTROL_ERROR U16
var RANGE_ERROR U16
var USAGE_ERROR U16
var ADMINISTRATIVE_EVENT U16
var STATISTIC_REPORT U16

func init() {
	SECURITY_ALARM = 1

	FAILURE_ALARM = 2

	DEGRADED_ALARM = 3

	EXCEPTION_ERROR = 4

	FLOW_CONTROL_ERROR = 5

	RANGE_ERROR = 6

	USAGE_ERROR = 7

	ADMINISTRATIVE_EVENT = 8

	STATISTIC_REPORT = 9

}
