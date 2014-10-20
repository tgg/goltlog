package goltlog

func NewLog(method string) (Log, error) {
	l := new(logger)
	if err := l.init(); err != nil {
		return nil, err
	}

	return l, nil
}
