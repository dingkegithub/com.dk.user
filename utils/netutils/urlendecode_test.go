package netutils

import (
	"net/url"
	"testing"
)

func TestStructToUrl(t *testing.T) {
	type Data struct {
		Name   string `json:"name"`
		Weight uint64 `json:"weight"`
		Enable bool   `json:"enable"`
	}

	type InnerData struct {
		A     string   `json:"a"`
		Data  Data     `json:"data"`
		Data1 *Data    `json:"data_1"`
		Arr   []string `json:"arr"`
	}

	data1 := &Data{
		Name:   "dkos",
		Weight: 100,
		Enable: false,
	}

	data2 := Data{
		Name:   "dkos",
		Weight: 100,
		Enable: false,
	}

	data3 := &InnerData{
		A:    "inner",
		Data: data2,
		Data1: &Data{
			Name:   "prtName",
			Weight: 0,
			Enable: true,
		},
		Arr: []string{"a", "b", "c"},
	}

	want := url.Values{}
	want.Add("name", "dkos")
	want.Add("weight", "100")
	want.Add("enable", "false")

	actual, err := StructToUrl(data1)
	if err != nil {
		t.Errorf("example error %s", err.Error())
		t.Fatal("invoke error")
	} else {
		t.Log("want is: ", want.Encode())
		t.Log("actual is: ", actual.Encode())

		if want.Encode() != actual.Encode() {
			t.Fatal("example failed")
		}
	}

	actual2, err := StructToUrl(data2)
	if err != nil {
		t.Errorf("example error %s", err.Error())
		t.Fatal("invoke error")
	} else {
		t.Log("want is: ", want.Encode())
		t.Log("actual is: ", actual2.Encode())

		if want.Encode() != actual2.Encode() {
			t.Fatal("example failed")
		}
	}
	wantInner := url.Values{}
	wantInner.Add("a", "inner")
	wantInner.Add("arr", "[\"a\",\"b\",\"c\"]")
	wantInner.Add("data", "{\"name\":\"dkos\",\"weight\":100,\"enable\":false}")
	wantInner.Add("data_1", "{\"name\":\"prtName\",\"weight\":0,\"enable\":true}")
	t.Log("want inner: ", wantInner.Encode())

	actual3, err := StructToUrl(data3)
	if err != nil {
		t.Errorf("example error %s", err.Error())
		t.Fatal("invoke error")
	} else {
		t.Log("want is: ", want.Encode())
		t.Log("actual is: ", actual3.Encode())

		if wantInner.Encode() != actual3.Encode() {
			t.Fatal("example failed")
		}
	}
}
