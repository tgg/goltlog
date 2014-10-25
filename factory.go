package goltlog

import (
	"errors"
)

func NewLog(method string) (Log, error) {
	switch (method) {
	case "local":
		l := new(logger)
		if err := l.init(); err != nil {
			return nil, err
		}

		return l, nil

	default:
		return nil, errors.New("Unsupported method")	
	}
}
