package goltlog

import (
	"time"
	"unsafe"
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
	lfa         LogFullAction
	as          AdministrativeState
	os          OperationalState
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

func (l *logger) GetNumRecords() (uint64, error) {
	if l.err != nil {
		return 0, l.err
	}
	return uint64(len(l.records)), nil
}

func (l *logger) GetLogFullAction() (LogFullAction, error) {
	if l.err != nil {
		return l.lfa, l.err
	}
	return l.lfa, nil
}

func (l *logger) GetAdministrativeState() (AdministrativeState, error) {
	if l.err != nil {
		return l.as, l.err
	}
	return l.as, nil
}

func (l *logger) GetOperationalState() (OperationalState, error) {
	if l.err != nil {
		return l.os, l.err
	}
	return l.os, nil
}

func (l *logger) GetAvailabilityStatus() (AvailabilityStatus, error) {
	st := AvailabilityStatus{
		OffDuty: l.as == LOCKED || l.os == DISABLED,
		LogFull: l.currentSize >= l.maxSize,
	}
	if l.err != nil {
		return st, l.err
	}
	return st, nil
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
		Info: *record,
	}
	l.records = append(l.records, lr)
	l.currentSize += uint64(unsafe.Sizeof(lr))
	return nil
}

func (l *logger) init() error {
	l.err = nil
	l.currentSize = 0
	l.maxSize = 0
	l.records = make([]LogRecord, 0, 10)
	l.lfa = HALT
	l.as = UNLOCKED
	l.os = ENABLED
	return nil
}
