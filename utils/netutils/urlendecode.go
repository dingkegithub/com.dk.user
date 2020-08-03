package netutils

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

var (
	ErrNotStruct = fmt.Errorf("not struct type")
)

func StructToUrl(s interface{}) (url.Values, error) {
	var sT reflect.Type
	var sV reflect.Value

	srcType := reflect.TypeOf(s)
	if srcType.Kind() == reflect.Ptr {
		sT = reflect.TypeOf(s).Elem()
		sV = reflect.ValueOf(s).Elem()
	} else {
		sT = reflect.TypeOf(s)
		sV = reflect.ValueOf(s)
	}

	if sT.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	var urlQuery url.Values
	urlQuery = make(map[string][]string)

	for i := 0; i < sT.NumField(); i++ {
		value := sV.Field(i).Interface()

		switch sT.Field(i).Type.Kind() {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Ptr:
			b, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			value = string(b)
		}

		qV := fmt.Sprintf("%v", value)
		if qV == "" {
			continue
		}

		tagName := sT.Field(i).Tag.Get("json")
		urlQuery.Add(tagName, qV)
	}

	return urlQuery, nil
}
