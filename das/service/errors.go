package service

import "errors"

var (
	ErrNotFound = errors.New("not found user")

	ErrParam = errors.New("parameter not valid")

	ErrExist = errors.New("record existed")

	ErrQuery = errors.New("query record failed")

	ErrSvcInner = errors.New("logic inner error")

	ErrOk = errors.New("process ok")
)

var (
	ErrCodeNotFound int64 = 40001

	ErrCodeParam int64 = 40002

	ErrCodeExist int64 = 40003

	ErrCodeQuery int64 = 40004

	ErrCodeSvcInner int64 = 50001

	ErrCodeOk int64 = 20000
)

func ErrMapToCode(err error) int64 {
	switch err {
	case ErrNotFound:
		return ErrCodeNotFound

	case ErrParam:
		return ErrCodeParam

	case ErrExist:
		return ErrCodeExist

	case ErrQuery:
		return ErrCodeQuery

	case ErrSvcInner:
		return ErrCodeSvcInner

	default:
		return ErrCodeOk
	}
}

func CodeMapToErr(code int64) error {
	switch code {
	case ErrCodeNotFound:
		return ErrNotFound

	case ErrCodeParam:
		return ErrParam

	case ErrCodeExist:
		return ErrExist

	case ErrCodeQuery:
		return ErrQuery

	case ErrCodeSvcInner:
		return ErrSvcInner

	default:
		return ErrOk
	}
}
