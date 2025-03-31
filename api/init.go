package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/ybbus/httpretry"
)

var client = httpretry.NewCustomClient(
	&http.Client{Timeout: 8 * time.Second},
	httpretry.WithMaxRetryCount(2),
	httpretry.WithRetryPolicy(func(statusCode int, err error) bool {
		return err != nil || statusCode >= 500 || statusCode == 0
	}),
	httpretry.WithBackoffPolicy(func(attemptNum int) time.Duration {
		return time.Duration(attemptNum+1) * time.Microsecond
	}),
)

func getRequest(urlstr string, data map[string]interface{}) (*http.Request, error) {
	u, err := url.Parse(urlstr)
	if err != nil {
		return nil, fmt.Errorf("getRequest  url.Parse  err:%+v", err)
	}

	var httpRequest *http.Request
	values := url.Values{}
	if data != nil {
		for k, v := range data {
			values.Add(k, fmt.Sprint(v))
		}
	}
	u.RawQuery = values.Encode()
	httpRequest, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("getRequest  http.NewRequest  err:%+v", err)
	}

	return httpRequest, nil
}

func postRequest(urlstr string, data map[string]interface{}) (*http.Request, error) {
	bytesData := []byte{}
	var err error
	var httpRequest *http.Request

	if data != nil {
		bytesData, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("postRequest  json.Marshal err:%+v", err)
		}
	}

	httpRequest, err = http.NewRequest("POST", urlstr, bytes.NewReader(bytesData))
	if err != nil {
		return nil, fmt.Errorf("postRequest  http.NewRequest err:%+v", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return httpRequest, nil
}

func parseResponse[T any](body []byte) (*JSONData[T], error) {
	data := &JSONData[T]{}
	err := json.Unmarshal(body, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func httpDo[T any](httpRequest *http.Request) (jsonData *JSONData[T], err error) {
	httpResponse, err := client.Do(httpRequest)
	if err != nil {
		httpResponse, err = client.Do(httpRequest)
		if err != nil {
			return nil, fmt.Errorf("httpDo 请求失败 err:%+v", err)
		}
	}
	defer httpResponse.Body.Close()

	body := []byte{}
	body, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("httpDo  ioutil.ReadAll err:%+v", err)
	}
	jsonData = &JSONData[T]{}

	jsonData, err = parseResponse[T](body)
	if err != nil {
		return nil, fmt.Errorf("httpDo  parseResponse err:%+v", err)
	}

	return jsonData, nil
}

func getResponse[T any](httpRequest *http.Request, s T) (jsonData *JSONData[T], err error) {
	jsonData, err = httpDo[T](httpRequest)
	if err != nil {
		err = fmt.Errorf("getResponse  httpDo  err:%+v", err)
		return
	}
	if jsonData.Code != 200 {
		return jsonData, fmt.Errorf("getResponse  jsonData.Code != 200  err:%+v", jsonData.Msg)
	}

	return
}
