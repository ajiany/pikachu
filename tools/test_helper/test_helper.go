package test_helper

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"time"
)

func JSONResp(resp *http.Response) *simplejson.Json {
	js, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	json, err := js.MarshalJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
	return js
}

func RunGetRequest(url string, query ...map[string]string) *http.Response {
	req := GetRequest(url, query...)

	return RunRequest(req)
}

func GetRequest(url string, query ...map[string]string) *http.Request {
	req, _ := http.NewRequest("GET", url, nil)

	q := req.URL.Query()
	for _, sq := range query {
		for k, v := range sq {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()

	return req
}

func RunRequest(req *http.Request) *http.Response {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	resp.Body.Close() //  must close
	fmt.Println(string(bodyBytes))
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return resp
}

func StringResp(resp *http.Response) string {
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(content)
}
