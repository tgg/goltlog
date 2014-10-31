package goltlog

import (
	"testing"
)

var (
	l Log
	err error
	nr uint64
)
func TestLocal(t *testing.T) {
	testImpl(t, "local")
}

func TestThrift(t *testing.T) {
	testImpl(t, "thrift")
}

func testImpl(t *testing.T, m string) {
	if l, err = NewLog(m); err != nil {
		t.Errorf("logger creation failed %v", err)
	}
	defer l.Destroy()
	r := &ProducerLogRecord{
		ProducerId: "0",
		ProducerName: "me",
		Level: FAILURE_ALARM,
		LogData: "Hello world!"}
	if err := l.WriteRecord(r); err != nil {
		t.Errorf("logger creation failed %v", err)
	}
	if nr, err = l.GetNRecords(); err != nil {
		t.Errorf("GetNRecords failed: %v", err)
	}
	if nr != 1 {
		t.Errorf("NumRecord: got %d, want %d", nr, 1)
	}
	if err := l.WriteRecord(r); err != nil {
		t.Errorf("logger creation failed %v", err)
	}
	if nr, err = l.GetNRecords(); err != nil {
		t.Errorf("GetNRecords failed: %v", err)
	}
	if nr != 2 {
		t.Errorf("NumRecord: got %d, want %d", nr, 2)
	}
}
