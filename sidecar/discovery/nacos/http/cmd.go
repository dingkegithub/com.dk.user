package http

import (
	"encoding/json"
	"fmt"
	"github.com/dingkegithub/com.dk.user/utils/netutils"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"net/http"
)

type Cmd struct {
	logger log.Logger
	nm *netutils.ClusterNodeManager
}

func NewCmd(logger log.Logger, nm *netutils.ClusterNodeManager) *Cmd {
	return &Cmd{
		logger: logger,
		nm:     nm,
	}
}

func (cmd *Cmd) Close() {
	cmd.nm.Close()
}

func (cmd *Cmd) buildUrl() string {
	nodeUrl, err := cmd.nm.Random()
	if err != nil {
		cmd.logger.Log("file", "cmd.go",
			"function", "buildUrl",
			"action", "random cluster node",
			"error", err)
		return ""
	}

	return fmt.Sprintf("http://%s", nodeUrl)
}


func (cmd *Cmd)CmdCreateInstance(req *RegisterInstanceRequest) error {
	url := cmd.buildUrl()
	if url == "" {
		return ErrNotFoundHealthyNode
	}

	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance?%s", url, queryValue.Encode())
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

func (cmd *Cmd)CmdDeleteInstance(req *DeregisterInstanceRequest) error {
	url := cmd.buildUrl()
	if url == "" {
		return ErrNotFoundHealthyNode
	}

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


func (cmd *Cmd)CmdUpdateInstance(req *ModifyInstanceRequest) error {
	url := cmd.buildUrl()
	if url == "" {
		return ErrNotFoundHealthyNode
	}

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

func (cmd *Cmd)CmdListInstance(req *ListInstanceRequest) (*ListInstanceResponse, error) {
	url := cmd.buildUrl()
	if url == "" {
		return nil, ErrNotFoundHealthyNode
	}

	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return nil, err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance/list?%s", url, queryValue.Encode())

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

func (cmd *Cmd)CmdDetailInstance(req *DetailInstanceRequest) (*DetailInstanceResponse, error) {
	url := cmd.buildUrl()
	if url == "" {
		return nil, ErrNotFoundHealthyNode
	}

	queryValue, err := netutils.StructToUrl(req)
	if err != nil {
		return nil, err
	}
	absUrl := fmt.Sprintf("%s/nacos/v1/ns/instance?%s", url, queryValue.Encode())

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

func (cmd *Cmd)CmdHeartbeatInstance(req *HeartbeatRequest) error {
	url := cmd.buildUrl()
	if url == "" {
		return ErrNotFoundHealthyNode
	}

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