package service

import "fmt"

var (
	ErrorCodeParaPassword int64 = 40000
	ErrorCodeSuccess int64 = 20000
	ErrorCodeServerInternal int64 = 50000
	ErrorCodeUnknown int64 = 60000
)

var (
	ErrorParaPassword = fmt.Errorf("password check failed")
	ErrorSuccess = fmt.Errorf("success")
	ErrorServerInternal = fmt.Errorf("server internal error")
	ErrorUnknown = fmt.Errorf("unknown error, please check server log")
)

func ErrCode(err error) int64 {
	switch err {
	case ErrorServerInternal:
		return ErrorCodeServerInternal

	case ErrorParaPassword:
		return ErrorCodeParaPassword

	case ErrorSuccess:
		return ErrorCodeSuccess

	default:
		return ErrorCodeUnknown

	}
}

func ErrMsg(code int64) error {
	switch code {
	case ErrorCodeServerInternal:
		return ErrorServerInternal

	case ErrorCodeParaPassword:
		return ErrorParaPassword

	case ErrorCodeSuccess:
		return ErrorSuccess

	default:
		return ErrorUnknown
	}
}
