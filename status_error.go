package statuserror

type StatusError interface {
	Unwrap() error
	StatusErr() *StatusErr
	Error() string
}

type StatusErrorWithServiceCode interface {
	ServiceCode() int
}
