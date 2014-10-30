package goltlog

import (
	"errors"
	"github.com/tgg/goltlog/rpc_thrift"
)

type thrift_cli struct {
	client *rpc_thrift.LogClient
}

func (p *thrift_cli) GetMaxSize() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.client.GetMaxSize(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return
}

func (p *thrift_cli) GetCurrentSize() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.client.GetCurrentSize(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return	
}

func (p *thrift_cli) GetNumRecords() (r uint64, err error) {
	var r_ rpc_thrift.U64
	if r_, err = p.client.GetNRecords(); err != nil {
		r = 0
		return
	}
	r = uint64(r_)
	return	
}

func (p *thrift_cli) GetLogFullAction() (r LogFullAction, err error) {
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

func (p *thrift_cli) GetAdministrativeState() (r AdministrativeState, err error) {
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

func (p *thrift_cli) GetAvailabilityStatus() (r AvailabilityStatus, err error) {
	var r_ *rpc_thrift.AvailabilityStatus
	if r_, err = p.client.GetAvailabilityStatus(); err != nil {
		r = AvailabilityStatus{OffDuty: true, LogFull: true}
		return
	}
	r = AvailabilityStatus{OffDuty: r_.OffDuty, LogFull: r_.LogFull}
	return
}

func (p *thrift_cli) GetOperationalState() (r OperationalState, err error) {
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

func (p *thrift_cli) SetMaxSize(s uint64) error {
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

func (p *thrift_cli) SetLogFullAction(a LogFullAction) error {
	var a_ rpc_thrift.LogFullAction
	if a == HALT {
		a_ = rpc_thrift.LogFullAction_HALT
	} else {
		a_ = rpc_thrift.LogFullAction_WRAP
	}
	return p.client.SetLogFullAction(a_)
}

func (p *thrift_cli) SetAdministrativeState(a AdministrativeState) error {
	var a_ rpc_thrift.AdministrativeState
	if a == LOCKED {
		a_ = rpc_thrift.AdministrativeState_locked
	} else {
		a_ = rpc_thrift.AdministrativeState_unlocked
	}
	return p.client.SetAdministrativeState(a_)
}

func (p *thrift_cli) ClearLog() error {
	return p.client.ClearLog()
}

func (p *thrift_cli) Destroy() (err error) {
	err = p.client.Destroy()
	if (err == nil) {
		err = p.client.Transport.Close()
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
	return p.client.WriteRecords(r_)
}

func (p *thrift_cli) WriteRecord(r *ProducerLogRecord) error {
	r_ := &rpc_thrift.ProducerLogRecord{
		ProducerId: r.ProducerId,
		ProducerName: r.ProducerName,
		Level: rpc_thrift.LogLevel(r.Level),
		LogData: r.LogData}
	return p.client.WriteRecord(r_)
}

func (p *thrift_cli) GetRecordIdFromTime(t LogTime) (RecordId, error) {
	t_ := &rpc_thrift.LogTime{
		Seconds: t.Seconds,
		Nanoseconds: t.Nanoseconds}
	i, err := p.client.GetRecordIdFromTime(t_)
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

func (p *thrift_cli) RetrieveRecords(i *RecordId, h *uint32) (r []LogRecord, err error) {
	i_ := rpc_thrift.RecordId(*i)
	h_ := rpc_thrift.U32(*h)
	r = nil
	var r_ *rpc_thrift.RecordIdU32LogRecordSequence
	if r_, err = p.client.RetrieveRecords(i_, h_); err != nil {
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
	if r_, err = p.client.RetrieveRecordsByLevel(i_, h_, l_); err != nil {
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
	if r_, err = p.client.RetrieveRecordsByProducerName(i_, h_, n); err != nil {
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
	if r_, err = p.client.RetrieveRecordsByProducerId(i_, h_, d); err != nil {
		return
	}
	*i = RecordId(r_.CurrentId)
	*h = uint32(r_.HowMany)
	r = convertLogRecord(r_.Result)
	return
}
