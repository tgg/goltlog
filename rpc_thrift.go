package goltlog

//go:generate thrift -out .. -gen go:thrift_import=git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift ltlog.thrift

import (
	"errors"
	"git-wip-us.apache.org/repos/asf/thrift.git/lib/go/thrift"
	"github.com/tgg/goltlog/rpc_thrift"
)

type thrift_cli struct {
	cli *rpc_thrift.LogClient
}

func (p *thrift_cli) GetMaxSize() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.cli.GetMaxSize(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return
}

func (p *thrift_cli) GetCurrentSize() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.cli.GetCurrentSize(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return	
}

func (p *thrift_cli) GetNRecords() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.cli.GetNRecords(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return	
}

func (p *thrift_cli) GetLogFullAction() (r LogFullAction, err error) {
	var r_ rpc_thrift.LogFullAction
	if r_, err = p.cli.GetLogFullAction(); err != nil {
		r = HALT
		return
	}
	if r_ == rpc_thrift.LogFullAction_HALT {
		r = HALT
	} else {
		r = WRAP
	}
	return
}

func (p *thrift_cli) GetAdministrativeState() (r AdministrativeState, err error) {
	var r_ rpc_thrift.AdministrativeState
	if r_, err = p.cli.GetAdministrativeState(); err != nil {
		r = LOCKED
		return
	}
	if r_ == rpc_thrift.AdministrativeState_locked {
		r = LOCKED
	} else {
		r = UNLOCKED
	}
	return
}

func (p *thrift_cli) GetAvailabilityStatus() (r AvailabilityStatus, err error) {
	var r_ *rpc_thrift.AvailabilityStatus
	if r_, err = p.cli.GetAvailabilityStatus(); err != nil {
		r = AvailabilityStatus{OffDuty: true, LogFull: true}
		return
	}
	r = AvailabilityStatus{OffDuty: r_.OffDuty, LogFull: r_.LogFull}
	return
}

func (p *thrift_cli) GetOperationalState() (r OperationalState, err error) {
	var r_ rpc_thrift.OperationalState
	if r_, err = p.cli.GetOperationalState(); err != nil {
		r = DISABLED
		return
	}
	if r_ == rpc_thrift.OperationalState_disabled {
		r = DISABLED
	} else {
		r = ENABLED
	}
	return
}

func (p *thrift_cli) SetMaxSize(s uint64) error {
	var ouch *rpc_thrift.InvalidParam
	var err error
	if ouch, err = p.cli.SetMaxSize(rpc_thrift.U64(s)); err != nil {
		return err
	}
	if ouch != nil {
		return errors.New(ouch.Details)
	}
	return nil
	
}

func (p *thrift_cli) SetLogFullAction(a LogFullAction) error {
	var a_ rpc_thrift.LogFullAction
	if a == HALT {
		a_ = rpc_thrift.LogFullAction_HALT
	} else {
		a_ = rpc_thrift.LogFullAction_WRAP
	}
	return p.cli.SetLogFullAction(a_)
}

func (p *thrift_cli) SetAdministrativeState(a AdministrativeState) error {
	var a_ rpc_thrift.AdministrativeState
	if a == LOCKED {
		a_ = rpc_thrift.AdministrativeState_locked
	} else {
		a_ = rpc_thrift.AdministrativeState_unlocked
	}
	return p.cli.SetAdministrativeState(a_)
}

func (p *thrift_cli) ClearLog() error {
	return p.cli.ClearLog()
}

func (p *thrift_cli) Destroy() (err error) {
	err = p.cli.Destroy()
	if (err == nil) {
		err = p.cli.Transport.Close()
	}
	return
}

func (p *thrift_cli) WriteRecords(r []ProducerLogRecord) error {
	r_ := make([]*rpc_thrift.ProducerLogRecord, len(r))
	for i, e := range(r) {
		r_[i] = &rpc_thrift.ProducerLogRecord{
			ProducerId: e.ProducerId,
			ProducerName: e.ProducerName,
			Level: rpc_thrift.LogLevel(e.Level),
			LogData: e.LogData}
	}
	return p.cli.WriteRecords(r_)
}

func (p *thrift_cli) WriteRecord(r *ProducerLogRecord) error {
	r_ := &rpc_thrift.ProducerLogRecord{
		ProducerId: r.ProducerId,
		ProducerName: r.ProducerName,
		Level: rpc_thrift.LogLevel(r.Level),
		LogData: r.LogData}
	return p.cli.WriteRecord(r_)
}

func (p *thrift_cli) GetRecordIdFromTime(t LogTime) (RecordId, error) {
	t_ := &rpc_thrift.LogTime{
		Seconds: t.Seconds,
		Nanoseconds: t.Nanoseconds}
	i, err := p.cli.GetRecordIdFromTime(t_)
	return RecordId(i), err
}

func convertLogRecord(t []*rpc_thrift.LogRecord) (r []LogRecord) {
	r = make ([]LogRecord, len(t))
	for j, e := range(t) {
		r[j].Id = RecordId(e.Id)
		r[j].Time.Seconds = e.Time.Seconds
		r[j].Time.Nanoseconds = e.Time.Nanoseconds
		r[j].Info.ProducerId = e.Info.ProducerId
		r[j].Info.ProducerName = e.Info.ProducerName
		r[j].Info.Level = LogLevel(e.Info.Level)
		r[j].Info.LogData = e.Info.LogData
	}
	return
}

func logRecordConvert(r []LogRecord) (t []*rpc_thrift.LogRecord)  {
	t = make([]*rpc_thrift.LogRecord, len(r))
	for i, e := range(r) {
		t[i] = &rpc_thrift.LogRecord{
			Id: rpc_thrift.RecordId(e.Id),
			Time: &rpc_thrift.LogTime{
				Seconds: e.Time.Seconds,
				Nanoseconds: e.Time.Nanoseconds},
			Info: &rpc_thrift.ProducerLogRecord{
				ProducerId: e.Info.ProducerId,
				ProducerName: e.Info.ProducerName,
				Level: rpc_thrift.LogLevel(e.Info.Level),
				LogData: e.Info.LogData}}
	}
	return
}

func (p *thrift_cli) RetrieveRecords(i *RecordId, h *uint32) (r []LogRecord, err error) {
	i_ := rpc_thrift.RecordId(*i)
	h_ := rpc_thrift.U32(*h)
	r = nil
	var r_ *rpc_thrift.RecordIdU32LogRecordSequence
	if r_, err = p.cli.RetrieveRecords(i_, h_); err != nil {
		return
	}
	*i = RecordId(r_.CurrentId)
	*h = uint32(r_.HowMany)
	r = convertLogRecord(r_.Result)
	return
}

func (p *thrift_cli) RetrieveRecordsByLevel(i *RecordId, h *uint32, l []LogLevel) (r []LogRecord, err error) {
	i_ := rpc_thrift.RecordId(*i)
	h_ := rpc_thrift.U32(*h)
	l_ := make([]rpc_thrift.LogLevel, len(l))
	for j, e := range (l) {
		l_[j] = rpc_thrift.LogLevel(e)
	}
	r = nil
	var r_ *rpc_thrift.RecordIdU32LogRecordSequence
	if r_, err = p.cli.RetrieveRecordsByLevel(i_, h_, l_); err != nil {
		return
	}
	*i = RecordId(r_.CurrentId)
	*h = uint32(r_.HowMany)
	r = convertLogRecord(r_.Result)
	return
}

func (p *thrift_cli) RetrieveRecordsByProducerName(i *RecordId, h *uint32, n []string) (r []LogRecord, err error) {
	i_ := rpc_thrift.RecordId(*i)
	h_ := rpc_thrift.U32(*h)
	r = nil
	var r_ *rpc_thrift.RecordIdU32LogRecordSequence
	if r_, err = p.cli.RetrieveRecordsByProducerName(i_, h_, n); err != nil {
		return
	}
	*i = RecordId(r_.CurrentId)
	*h = uint32(r_.HowMany)
	r = convertLogRecord(r_.Result)
	return
}

func (p *thrift_cli) RetrieveRecordsByProducerId(i *RecordId, h *uint32, d []string) (r []LogRecord, err error) {
	i_ := rpc_thrift.RecordId(*i)
	h_ := rpc_thrift.U32(*h)
	r = nil
	var r_ *rpc_thrift.RecordIdU32LogRecordSequence
	if r_, err = p.cli.RetrieveRecordsByProducerId(i_, h_, d); err != nil {
		return
	}
	*i = RecordId(r_.CurrentId)
	*h = uint32(r_.HowMany)
	r = convertLogRecord(r_.Result)
	return
}


// Handler to a real implementation.
type thrift_srv struct {
	h   Log
	srv thrift.TServer
}

func (p *thrift_srv) GetMaxSize() (r rpc_thrift.U64, err error) {
	var r_ uint64
	r_, err = p.h.GetMaxSize()
	if err == nil {
		r = rpc_thrift.U64(r_)
	}
	return
}

func (p *thrift_srv) GetCurrentSize() (r rpc_thrift.U64, err error) {
	var r_ uint64
	r_, err = p.h.GetCurrentSize()
	if err == nil {
		r = rpc_thrift.U64(r_)
	}
	return
}

func (p *thrift_srv) GetNRecords() (r rpc_thrift.U64, err error) {
	var r_ uint64
	r_, err = p.h.GetNRecords()
	if err == nil {
		r = rpc_thrift.U64(r_)
	}
	return	
}

func (p *thrift_srv) GetLogFullAction() (r rpc_thrift.LogFullAction, err error) {
	var r_ LogFullAction
	r_, err = p.h.GetLogFullAction()
	if err == nil {
		if r_ == HALT {
			r = rpc_thrift.LogFullAction_HALT
		} else {
			r = rpc_thrift.LogFullAction_WRAP
		}
	}
	return
}

func (p *thrift_srv) GetAvailabilityStatus() (r *rpc_thrift.AvailabilityStatus, err error) {
	var r_ AvailabilityStatus
	r_, err = p.h.GetAvailabilityStatus()
	if err == nil {
		r = &rpc_thrift.AvailabilityStatus{
			OffDuty: r_.OffDuty,
			LogFull: r_.LogFull}
	}
	return
}

func (p *thrift_srv) GetAdministrativeState() (r rpc_thrift.AdministrativeState, err error) {
	r_, err := p.h.GetAdministrativeState()
	if err == nil {
		if r_ == LOCKED {
			r = rpc_thrift.AdministrativeState_locked
		} else {
			r = rpc_thrift.AdministrativeState_unlocked
		}
	}
	return
}

func (p *thrift_srv) GetOperationalState() (r rpc_thrift.OperationalState, err error) {
	var r_ OperationalState
	r_, err = p.h.GetOperationalState()
	if err == nil {
		if r_ == DISABLED {
			r = rpc_thrift.OperationalState_disabled
		} else {
			r = rpc_thrift.OperationalState_enabled
		}
	}
	return
}

func (p *thrift_srv) SetMaxSize(size rpc_thrift.U64) (ouch *rpc_thrift.InvalidParam, err error) {
	err = p.h.SetMaxSize(uint64(size))
	if err != nil {
		ouch = &rpc_thrift.InvalidParam{Details: err.Error()}
	}
	return
}

func (p *thrift_srv) SetLogFullAction(action rpc_thrift.LogFullAction) (err error) {
	var a_ LogFullAction
	if action == rpc_thrift.LogFullAction_HALT {
		a_ = HALT
	} else {
		a_ = WRAP
	}
	return p.h.SetLogFullAction(a_)
}

func (p *thrift_srv) SetAdministrativeState(state rpc_thrift.AdministrativeState) (err error) {
	var s_ AdministrativeState
	if state == rpc_thrift.AdministrativeState_unlocked {
		s_ = UNLOCKED
	} else {
		s_ = LOCKED
	}
	return p.h.SetAdministrativeState(s_)
}

func (p *thrift_srv) ClearLog() error {
	return p.h.ClearLog()
}

func (p *thrift_srv) Destroy() (err error) {
	if err = p.h.Destroy(); err == nil {
		if (p.srv != nil) {
			p.srv.Stop()
		}
	}
	return
}

func (p *thrift_srv) GetRecordIdFromTime(fromTime *rpc_thrift.LogTime) (r rpc_thrift.RecordId, err error) {
	t := LogTime{
		Seconds: int32(fromTime.Seconds),
		Nanoseconds: int32(fromTime.Nanoseconds)}
	if r_, err := p.h.GetRecordIdFromTime(t); err == nil {
		r = rpc_thrift.RecordId(r_)
	}
	return
}

func (p *thrift_srv) RetrieveRecords(currentId rpc_thrift.RecordId, howMany rpc_thrift.U32) (r *rpc_thrift.RecordIdU32LogRecordSequence, err error) {
	i_ := new(RecordId)
	*i_ = RecordId(currentId)
	h_ := new(uint32)
	*h_ = uint32(howMany)
	if r_, err := p.h.RetrieveRecords(i_, h_); err == nil {
		r = &rpc_thrift.RecordIdU32LogRecordSequence{
			CurrentId: rpc_thrift.RecordId(*i_),
			HowMany: rpc_thrift.U32(*h_),
			Result: logRecordConvert(r_)}
	}
	return
}

func (p *thrift_srv) RetrieveRecordsByLevel(currentId rpc_thrift.RecordId, howMany rpc_thrift.U32, valueList rpc_thrift.LogLevelSequence) (r *rpc_thrift.RecordIdU32LogRecordSequence, err error) {
	i_ := new(RecordId)
	*i_ = RecordId(currentId)
	h_ := new(uint32)
	*h_ = uint32(howMany)
	s_ := make([]LogLevel, len(valueList))
	for j, e := range(valueList) {
		s_[j] = LogLevel(e)
	}
	if r_, err := p.h.RetrieveRecordsByLevel(i_, h_, s_); err == nil {
		r = &rpc_thrift.RecordIdU32LogRecordSequence{
			CurrentId: rpc_thrift.RecordId(*i_),
			HowMany: rpc_thrift.U32(*h_),
			Result: logRecordConvert(r_)}
	}
	return
}

func (p *thrift_srv) RetrieveRecordsByProducerId(currentId rpc_thrift.RecordId, howMany rpc_thrift.U32, valueList rpc_thrift.StringSeq) (r *rpc_thrift.RecordIdU32LogRecordSequence, err error) {
	i_ := new(RecordId)
	*i_ = RecordId(currentId)
	h_ := new(uint32)
	*h_ = uint32(howMany)
	if r_, err := p.h.RetrieveRecordsByProducerId(i_, h_, valueList); err == nil {
		r = &rpc_thrift.RecordIdU32LogRecordSequence{
			CurrentId: rpc_thrift.RecordId(*i_),
			HowMany: rpc_thrift.U32(*h_),
			Result: logRecordConvert(r_)}
	}
	return
}

func (p *thrift_srv) RetrieveRecordsByProducerName(currentId rpc_thrift.RecordId, howMany rpc_thrift.U32, valueList rpc_thrift.StringSeq) (r *rpc_thrift.RecordIdU32LogRecordSequence, err error) {
	i_ := new(RecordId)
	*i_ = RecordId(currentId)
	h_ := new(uint32)
	*h_ = uint32(howMany)
	if r_, err := p.h.RetrieveRecordsByProducerName(i_, h_, valueList); err == nil {
		r = &rpc_thrift.RecordIdU32LogRecordSequence{
			CurrentId: rpc_thrift.RecordId(*i_),
			HowMany: rpc_thrift.U32(*h_),
			Result: logRecordConvert(r_)}
	}
	return
}

func (p *thrift_srv) WriteRecords(records rpc_thrift.ProducerLogRecordSequence) (err error) {
	r_ := make([]ProducerLogRecord, len(records))
	for i, e := range(records) {
		r_[i].ProducerId = e.ProducerId
		r_[i].ProducerName = e.ProducerName
		r_[i].Level = LogLevel(e.Level)
		r_[i].LogData = e.LogData}
	return p.h.WriteRecords(r_)
}

func (p *thrift_srv) WriteRecord(record *rpc_thrift.ProducerLogRecord) (err error) {
	return p.h.WriteRecord(&ProducerLogRecord{
		ProducerId: record.ProducerId,
		ProducerName: record.ProducerName,
		Level: LogLevel(record.Level),
		LogData: record.LogData})
}
