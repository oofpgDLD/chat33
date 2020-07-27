package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if nil != resp {
		defer resp.Body.Close()
	}

	if nil != err {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func JSONRPCPost(url string, payload []byte) ([]byte, error) {
	response, err := http.Post(url, "text/json", bytes.NewBuffer(payload))
	if nil != response {
		defer response.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http request error: %v", response.StatusCode)
	}

	return ioutil.ReadAll(response.Body)
}
