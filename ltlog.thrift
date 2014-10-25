namespace go gotltlog

// No unsigned types in thrift. Recommendations are:
// 1. use next bigger type; or
// 2. cast between signed and unsigned
//
// We pick option 2. here and define false types:
typedef i16 u16
typedef i32 u32
typedef i64 u64

const u16 SECURITY_ALARM = 1
const u16 FAILURE_ALARM = 2
const u16 DEGRADED_ALARM = 3
const u16 EXCEPTION_ERROR = 4
const u16 FLOW_CONTROL_ERROR = 5
const u16 RANGE_ERROR = 6
const u16 USAGE_ERROR = 7
const u16 ADMINISTRATIVE_EVENT = 8
const u16 STATISTIC_REPORT = 9

typedef u16 LogLevel
enum OperationalState {disabled, enabled}
enum AdministrativeState {locked, unlocked}
enum LogFullAction {WRAP, HALT}

typedef u64 RecordId
struct LogTime {
  1: i32 seconds
  2: i32 nanoseconds
}
struct AvailabilityStatus {
  1: bool off_duty
  2: bool log_full
}
struct ProducerLogRecord {
  1: string producerId
  2: string producerName
  3: LogLevel level
  4: string logData
}
struct LogRecord {
  1: RecordId id
  2: LogTime time
  3: ProducerLogRecord info
}
typedef list<LogRecord> LogRecordSequence
typedef list<ProducerLogRecord> ProducerLogRecordSequence
typedef list<LogLevel> LogLevelSequence
typedef list<string> StringSeq
exception InvalidParam {
  1: string details
}
service LogStatus {
  u64 get_max_size()
  u64 get_current_size()
  u64 get_n_records()
  LogFullAction get_log_full_action()
  AvailabilityStatus get_availability_status()
  AdministrativeState get_administrative_state()
  OperationalState get_operational_state()
}
service LogAdministrator extends LogStatus {
  void set_max_size(1:u64 size) throws (1:InvalidParam ouch)
  void set_log_full_action(1:LogFullAction action)
  void set_administrative_state(1:AdministrativeState state)
  void clear_log()
  void destroy ()
}
// No out or inout parameters in Thrift. Hence we use a structure.
struct RecordIdU32LogRecordSequence {
  1: RecordId currentId
  2: u32 howMany
  3: LogRecordSequence result
}
service LogConsumer extends LogStatus {
  RecordId get_record_id_from_time(1:LogTime fromTime)
  RecordIdU32LogRecordSequence retrieve_records(1:RecordId currentId, 2:u32 howMany)
  RecordIdU32LogRecordSequence retrieve_records_by_level(1:RecordId currentId, 2:u32 howMany, 3:LogLevelSequence valueList)
  RecordIdU32LogRecordSequence retrieve_records_by_producer_id(1:RecordId currentId, 2:u32 howMany, 3:StringSeq valueList)
  RecordIdU32LogRecordSequence retrieve_records_by_producer_name(1:RecordId currentId, 2:u32 howMany, 3:StringSeq valueList)
}
service LogProducer extends LogStatus {
  oneway void write_records(1:ProducerLogRecordSequence records)
  oneway void write_record(1:ProducerLogRecord record)
}
// extends limited to a single type, so we copy/paste others here
service Log extends LogConsumer {
  void set_max_size(1:u64 size) throws (1:InvalidParam ouch)
  void set_log_full_action(1:LogFullAction action)
  void set_administrative_state(1:AdministrativeState state)
  void clear_log()
  void destroy ()
  oneway void write_records(1:ProducerLogRecordSequence records)
  oneway void write_record(1:ProducerLogRecord record)
}
