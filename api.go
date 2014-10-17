package goltlog

type LogLevel uint16
type RecordId uint64

const (
	SECURITY_ALARM       LogLevel = 1 + iota
	FAILURE_ALARM
	DEGRADED_ALARM
	EXCEPTION_ERROR
	FLOW_CONTROL_ERROR
	RANGE_ERROR
	USAGE_ERROR
	ADMINISTRATIVE_EVENT
	STATISTIC_REPORT
)

type LogTime struct {
	Seconds     int64
	NanoSeconds int64
}

type ProducerLogRecord struct {
	ProducerId   string
	ProducerName string
	Level        LogLevel
	LogData      string
}

type LogRecord struct {
	Id   RecordId
	Time LogTime
	Info ProducerLogRecord
}

type LogStatus interface {
	GetMaxSize() (uint64, error)
	GetCurrentSize() (uint64, error)
	GetNRecords() (uint64, error)
	// TODO: Following fields need enum
}

// TODO: see errors again... how can you guess type
// if errors are serialized as string?
type errInvalidParam struct {
	error
	Details string
}

type LogAdministrator interface {
	LogStatus
	SetMaxSize(uint64) error
	ClearLog() error
	Destroy() error
}

type LogProducer interface {
	LogStatus
	WriteRecords([]ProducerLogRecord) error
	WriteRecord(*ProducerLogRecord) error
}
