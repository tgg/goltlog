package goltlog

func NewLogger(method string) (LogProducer, error) {
	l := new(logger)
	if err := l.init(); err != nil {
		return nil, err
	}

	return l, nil
}
