package goltlog

type LogLevel uint16
type RecordId uint64
type OperationalState bool
type AdministrativeState bool
type LogFullAction bool

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

const (
	DISABLED OperationalState = false
	ENABLED                   = true
)

const (
	LOCKED AdministrativeState = false
	UNLOCKED                   = true
)

const (
	HALT LogFullAction = false
	WRAP               = true
)

type AvailabilityStatus struct {
	OffDuty bool
	LogFull bool
}

type LogTime struct {
	Seconds     int32
	Nanoseconds int32
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
	GetNumRecords() (uint64, error)
	GetLogFullAction() (LogFullAction, error)
	GetAdministrativeState() (AdministrativeState, error)
	GetAvailabilityStatus() (AvailabilityStatus, error)
	GetOperationalState() (OperationalState, error)
}

type LogAdministrator interface {
	LogStatus
	SetMaxSize(uint64) error
	SetLogFullAction(LogFullAction) error
	SetAdministrativeState(AdministrativeState) error
	ClearLog() error
	Destroy() error
}

type LogProducer interface {
	LogStatus
	WriteRecords([]ProducerLogRecord) error
	WriteRecord(*ProducerLogRecord) error
}

type LogConsumer interface {
	LogStatus
	GetRecordIdFromTime(LogTime) (RecordId, error)
	RetrieveRecords(*RecordId, *uint32) ([]LogRecord, error)
	RetrieveRecordsByLevel(*RecordId, *uint32, []LogLevel) ([]LogRecord, error)
	RetrieveRecordsByProducerName(*RecordId, *uint32, []string) ([]LogRecord, error)
	RetrieveRecordsByProducerId(*RecordId, *uint32, []string) ([]LogRecord, error)
}

// Does not compile because of duplicate method (diamond pattern on LogStatus)
//
// type Log interface {
// 	LogProducer
// 	LogAdministrator
// 	LogConsumer
// }
//
//
// This means that a naive but working IDL to Go compiler should copy base
// method in interface declaration rather than embed

type Log interface {
	GetMaxSize() (uint64, error)
	GetCurrentSize() (uint64, error)
	GetNumRecords() (uint64, error)
	GetLogFullAction() (LogFullAction, error)
	GetAdministrativeState() (AdministrativeState, error)
	GetAvailabilityStatus() (AvailabilityStatus, error)
	GetOperationalState() (OperationalState, error)
	SetMaxSize(uint64) error
	SetLogFullAction(LogFullAction) error
	SetAdministrativeState(AdministrativeState) error
	ClearLog() error
	Destroy() error
	WriteRecords([]ProducerLogRecord) error
	WriteRecord(*ProducerLogRecord) error
	GetRecordIdFromTime(LogTime) (RecordId, error)
	RetrieveRecords(*RecordId, *uint32) ([]LogRecord, error)
	RetrieveRecordsByLevel(*RecordId, *uint32, []LogLevel) ([]LogRecord, error)
	RetrieveRecordsByProducerName(*RecordId, *uint32, []string) ([]LogRecord, error)
	RetrieveRecordsByProducerId(*RecordId, *uint32, []string) ([]LogRecord, error)
}
