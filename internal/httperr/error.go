package httperr

type Error struct {
	error
	code int
}

func New(code int, err error) error {
	return &Error{
		error: err,
		code:  code,
	}
}

func (e *Error) StatusCode() int {
	return e.code
}

func (e *Error) Unwrap() error {
	return e.error
}
