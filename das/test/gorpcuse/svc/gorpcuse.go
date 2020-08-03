package svc

import "strings"

type StringRequest struct {
	A string
	B string
}

type Service interface {
	Concat(request StringRequest, ret *string) error

	Diff(request StringRequest, ret *string) error
}

type StringService struct {

}

func (s StringService) Concat(request StringRequest, res *string) error {
	*res = request.A + request.B
	return nil
}

func (s StringService) Diff(request StringRequest, res *string) error {
	ret := ""

	if len(request.A) >= len(request.B) {
		for _, char := range request.B {
			if strings.Contains(request.A, string(char)) {
				ret = ret + string(char)
			}
		}
	} else {
		for _, char := range request.A {
			if strings.Contains(request.B, string(char)) {
				ret = ret + string(char)
			}
		}
	}

	*res = ret
	return nil
}

