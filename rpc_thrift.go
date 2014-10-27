package goltlog

import (
	"errors"
	"github.com/tgg/goltlog/rpc_thrift"
)

type thrift_logger struct {
	client *rpc_thrift.LogClient
}

func (p *thrift_logger) GetMaxSize() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.client.GetMaxSize(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return
}

func (p *thrift_logger) GetCurrentSize() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.client.GetCurrentSize(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return	
}

func (p *thrift_logger) GetNumRecords() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.client.GetNRecords(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return	
}

func (p *thrift_logger) GetLogFullAction() (r LogFullAction, err error) {
	var r_ rpc_thrift.LogFullAction
	if r_, err = p.client.GetLogFullAction(); err != nil {
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

func (p *thrift_logger) GetAdministrativeState() (r AdministrativeState, err error) {
	var r_ rpc_thrift.AdministrativeState
	if r_, err = p.client.GetAdministrativeState(); err != nil {
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

func (p *thrift_logger) GetAvailabilityStatus() (r AvailabilityStatus, err error) {
	var r_ *rpc_thrift.AvailabilityStatus
	if r_, err = p.client.GetAvailabilityStatus(); err != nil {
		r = AvailabilityStatus{OffDuty: true, LogFull: true}
		return
	}
	r = AvailabilityStatus{OffDuty: r_.OffDuty, LogFull: r_.LogFull}
	return
}

func (p *thrift_logger) GetOperationalState() (r OperationalState, err error) {
	var r_ rpc_thrift.OperationalState
	if r_, err = p.client.GetOperationalState(); err != nil {
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

func (p *thrift_logger) SetMaxSize(s uint64) error {
	var ouch *rpc_thrift.InvalidParam
	var err error
	if ouch, err = p.client.SetMaxSize(rpc_thrift.U64(s)); err != nil {
		return err
	}
	if ouch != nil {
		return errors.New(ouch.Details)
	}
	return nil
	
}

func (p *thrift_logger) SetLogFullAction(a LogFullAction) error {
	var a_ rpc_thrift.LogFullAction
	if a == HALT {
		a_ = rpc_thrift.LogFullAction_HALT
	} else {
		a_ = rpc_thrift.LogFullAction_WRAP
	}
	return p.client.SetLogFullAction(a_)
}

func (p *thrift_logger) SetAdministrativeState(a AdministrativeState) error {
	var a_ rpc_thrift.AdministrativeState
	if a == LOCKED {
		a_ = rpc_thrift.AdministrativeState_locked
	} else {
		a_ = rpc_thrift.AdministrativeState_unlocked
	}
	return p.client.SetAdministrativeState(a_)
}

func (p *thrift_logger) ClearLog() error {
	return p.client.ClearLog()
}

func (p *thrift_logger) Destroy() error {
	return p.client.Destroy()
}

func (p *thrift_logger) WriteRecords(r []ProducerLogRecord) error {
	r_ := make([]*rpc_thrift.ProducerLogRecord, len(r))
	for i, e := range(r) {
		r_[i] = &rpc_thrift.ProducerLogRecord{
			ProducerId: e.ProducerId,
			ProducerName: e.ProducerName,
			Level: rpc_thrift.LogLevel(e.Level),
			LogData: e.LogData}
	}
	return p.client.WriteRecords(r_)
}

func (p *thrift_logger) WriteRecord(r *ProducerLogRecord) error {
	r_ := &rpc_thrift.ProducerLogRecord{
		ProducerId: r.ProducerId,
		ProducerName: r.ProducerName,
		Level: rpc_thrift.LogLevel(r.Level),
		LogData: r.LogData}
	return p.client.WriteRecord(r_)
}

func (p *thrift_logger) GetRecordIdFromTime(t LogTime) (RecordId, error) {
	t_ := &rpc_thrift.LogTime{
		Seconds: t.Seconds,
		Nanoseconds: t.Nanoseconds}
	i, err := p.client.GetRecordIdFromTime(t_)
	return RecordId(i), err
}

func (p *thrift_logger) RetrieveRecords(i *RecordId, h *uint32) (r []LogRecord, err error) {
	i_ := rpc_thrift.RecordId(*i)
	h_ := rpc_thrift.U32(*h)
	r = nil
	var r_ *rpc_thrift.RecordIdU32LogRecordSequence
	if r_, err = p.client.RetrieveRecords(i_, h_); err != nil {
		return
	}
	*i = RecordId(r_.CurrentId)
	*h = uint32(r_.HowMany)
	r = make ([]LogRecord, len(r_.Result))
	for j, e := range(r_.Result) {
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

func (p *thrift_logger) RetrieveRecordsByLevel(*RecordId, *uint32, []LogLevel) ([]LogRecord, error) {
	return nil, nil
}

func (p *thrift_logger) RetrieveRecordsByProducerName(*RecordId, *uint32, []string) ([]LogRecord, error) {
	return nil, nil
}

func (p *thrift_logger) RetrieveRecordsByProducerId(*RecordId, *uint32, []string) ([]LogRecord, error) {
	return nil, nil
}
