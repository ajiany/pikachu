package json

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"net/http"
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
