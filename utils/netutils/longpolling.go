package netutils

import (
	"io/ioutil"
	"net/http"
	"time"
)

func LongPolling(watchUrl string, timeout time.Duration) ([]byte, error) {

	for {
		client := http.Client{Timeout: timeout}
		resp, err := client.Get(watchUrl)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 304 {
			_ = resp.Body.Close()
			continue
		}

		if resp.StatusCode != 200 {
			_ = resp.Body.Close()
			time.Sleep(time.Second * 1)
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return data, err
	}
}
