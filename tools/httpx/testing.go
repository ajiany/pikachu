package httpx

import (
	"net/http/httptest"
)

func NewTestResponse(code int, body []byte) *Response {
	r := httptest.NewRecorder()
	r.Write(body)
	r.Flush()
	r.Code = code

	return newResponse(r.Result())
}
