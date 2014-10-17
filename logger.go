package goltlog

import (
	"time"
)

func currentTime() LogTime {
	now := time.Now().UTC()
	return LogTime{
		Seconds: now.Unix(),
		NanoSeconds: int64(now.Nanosecond()),
	}
}

type logger struct {
	err         error
	currentSize uint64
	maxSize     uint64
	records     []LogRecord
}

func (l *logger) GetMaxSize() (uint64, error) {
	if l.err != nil {
		return 0, l.err
	}
	return l.maxSize, nil
}

func (l *logger) GetCurrentSize() (uint64, error) {
	if l.err != nil {
		return 0, l.err
	}
	return l.currentSize, nil
}

func (l *logger) GetNRecords() (uint64, error) {
	if l.err != nil {
		return 0, l.err
	}
	return uint64(len(l.records)), nil
}

func (l *logger) SetMaxSize(s uint64) error {
	if l.err != nil {
		return l.err
	}
	l.maxSize = s
	return nil
}

func (l *logger) ClearLog() error {
	if l.err != nil {
		return l.err
	}
	l.init()
	return nil
}

func (l *logger) Destroy() error {
	if err := l.ClearLog(); err != nil {
		return err
	}
	return nil
}

func (l *logger) WriteRecords(records []ProducerLogRecord) error {
	if l.err != nil {
		return l.err
	}
	for _, r := range records {
		if err := l.WriteRecord(&r); err != nil {
			return err
		}
	}
	return nil
}

func (l *logger) WriteRecord(record *ProducerLogRecord) error {
	if l.err != nil {
		return l.err
	}
	lr := LogRecord{
		Id:   RecordId(len(l.records) + 1),
		Time: currentTime(),
		Info: *record}
	l.records = append(l.records, lr)
	return nil
}

func (l *logger) init() error {
	l.err = nil
	l.currentSize = 0
	l.maxSize = 0
	l.records = make([]LogRecord, 0, 10)
	return nil
}
