package goltlog

import (
	"testing"
)

var (
	l Log
	err error
	nr uint64
)

func TestApi(t *testing.T) {
	if l, err = NewLog("local"); err != nil {
		t.Errorf("logger creation failed %v", err)
	}
	r := &ProducerLogRecord{
		ProducerId: "0",
		ProducerName: "me",
		Level: FAILURE_ALARM,
		LogData: "Hello world!"}
	if err := l.WriteRecord(r); err != nil {
		t.Errorf("logger creation failed %v", err)
	}
	if nr, err = l.GetNumRecords(); err != nil {
		t.Errorf("GetNumRecords failed %v", err)
	}
	if nr != 1 {
		t.Errorf("NumRecord: got %d, want %d", nr, 1)
	}
	if err := l.WriteRecord(r); err != nil {
		t.Errorf("logger creation failed %v", err)
	}
	if nr, err = l.GetNumRecords(); err != nil {
		t.Errorf("GetNumRecords failed %v", err)
	}
	if nr != 2 {
		t.Errorf("NumRecord: got %d, want %d", nr, 2)
	}
}
