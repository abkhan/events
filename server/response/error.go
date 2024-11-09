package response

import "errors"

type Error struct {
	Message string // user
	Status  string // short descr
	Code    int

	data interface{}
	err  error // log
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) copy() *Error {
	return &Error{
		Message: e.Message,
		Status:  e.Status,
		Code:    e.Code,
		data:    e.data,
		err:     e.err,
	}
}

func (e *Error) WithError(err error) *Error {
	n := e.copy()
	n.err = err
	return n
}

func (e *Error) WithData(data interface{}) *Error {
	n := e.copy()
	n.data = data
	return n
}

func (e *Error) WithMessage(msg string) *Error {
	n := e.copy()
	n.Message = msg
	return n
}

func (e *Error) fillMissedFields() {
	if e.err == nil {
		e.err = errors.New(e.Message)
	}
}
