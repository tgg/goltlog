package goltlog

import (
	"time"
	"unsafe"
)

var (
	empty = make([]LogRecord, 0, 0)
)

func min(a, b int) int {
   if a < b {
      return a
   }
   return b
}

func currentTime() LogTime {
	now := time.Now().UTC()
	return LogTime{
		Seconds:     now.Unix(),
		NanoSeconds: int64(now.Nanosecond()),
	}
}

func (t LogTime) geq(o LogTime) bool {
	return t.Seconds > o.Seconds || (t.Seconds == o.Seconds && t.NanoSeconds >= o.NanoSeconds)
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

func (l *logger) GetRecordIdFromTime(t LogTime) (RecordId, error) {
	if l.err != nil {
		return RecordId(0), l.err
	}
	id := RecordId(len(l.records) + 1)
	for _, r := range l.records {
		if r.Time.geq(t) {
			id = r.Id
			break
		}
	}
	return id, nil
}

func (l *logger) RetrieveRecords(currentId *RecordId, howMany *uint32) ([]LogRecord, error) {
	if l.err != nil {
		return nil, l.err
	}
	if uint64(*currentId) >= uint64(len(l.records)) {
		return empty, nil
	}
	u := min(int(*currentId) + int(*howMany), len(l.records) - int(*currentId))
	// Security risk: we return underlying structure.
	r := l.records[int(*currentId):u]
	*howMany = uint32(u - int(*currentId))
	*currentId = RecordId(u + 1)
	return r, nil
}

func (l *logger) RetrieveRecordsByLevel(currentId *RecordId, howMany *uint32, levels []LogLevel) ([]LogRecord, error) {
	// TODO real implementation
	return l.RetrieveRecords(currentId, howMany)
}

func (l *logger) RetrieveRecordsByProducerName(currentId *RecordId, howMany *uint32, names []string) ([]LogRecord, error) {
	// TODO real implementation
	return l.RetrieveRecords(currentId, howMany)
}

func (l *logger) RetrieveRecordsByProducerId(currentId *RecordId, howMany *uint32, ids []string) ([]LogRecord, error) {
	// TODO real implementation
	return l.RetrieveRecords(currentId, howMany)
}

func (l *logger) SetMaxSize(s uint64) error {
	if l.err != nil {
		return l.err
	}
	l.maxSize = s
	return nil
}

func (l *logger) SetLogFullAction(a LogFullAction) error {
	if l.err != nil {
		return l.err
	}
	l.lfa = a
	return nil
}

func (l *logger) SetAdministrativeState(s AdministrativeState) error {
	if l.err != nil {
		return l.err
	}
	l.as = s
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
