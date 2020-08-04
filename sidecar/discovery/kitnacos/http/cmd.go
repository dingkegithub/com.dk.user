package http

import (
	"encoding/json"
	"fmt"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"io/ioutil"
	"net/http"
)

func CmdCreateInstance(url string, req *RegisterInstanceRequest) error {
	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance?%s", url, queryValue.Encode())
	fmt.Println("cmd create instance url: ", absUrl)
	resp, err := http.DefaultClient.Post(absUrl, "application/x-www-form-urlencode", nil)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) == "ok" {
		return nil
	}

	return fmt.Errorf("%s", string(body))
}
func CmdDeleteInstance(url string, req *DeregisterInstanceRequest) error {
	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance?%s", url, queryValue.Encode())

	request, err := http.NewRequest("DELETE", absUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencode")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) == "ok" {
		return nil
	}

	return fmt.Errorf("%s", string(body))
}


func CmdUpdateInstance(url string, req *ModifyInstanceRequest) error {
	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance?%s", url, queryValue.Encode())

	request, err := http.NewRequest("PUT", absUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencode")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) == "ok" {
		return nil
	}

	return fmt.Errorf("%s", string(body))
}

func CmdListInstance(url string, req *ListInstanceRequest) (*ListInstanceResponse, error) {
	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return nil, err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance/list?%s", url, queryValue.Encode())
	fmt.Println("URL: ", absUrl)

	request, err := http.NewRequest("GET", absUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencode")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respList := &ListInstanceResponse{}
	err = json.Unmarshal(body, &respList)
	if err != nil {
		return nil, err
	}

	return respList, nil
}

func CmdDetailInstance(url string, req *DetailInstanceRequest) (*DetailInstanceResponse, error) {
	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return nil, err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance?%s", url, queryValue.Encode())
	fmt.Println("URL: ", absUrl)

	request, err := http.NewRequest("GET", absUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencode")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	detail := &DetailInstanceResponse{}
	err = json.Unmarshal(body, &detail)
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func CmdHeartbeatInstance(url string, req *HeartbeatRequest) error {
	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance/beat?%s", url, queryValue.Encode())

	request, err := http.NewRequest("PUT", absUrl, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencode")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) == "ok" {
		return nil
	}

	return fmt.Errorf("%s", string(body))
}

