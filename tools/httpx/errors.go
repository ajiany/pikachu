package httpx

import "fmt"

type RequestFailErr struct {
	Resp *Response
}

func NewRequestFailErr(r *Response) RequestFailErr {
	return RequestFailErr{Resp: r}
}

func IsRequestFailErr(err error) bool {
	_, ok := err.(RequestFailErr)
	return ok
}

func (err RequestFailErr) Error() string {
	body, _ := err.Resp.ParsedString()

	return fmt.Sprintf("code: %d, msg: %s", err.Resp.Resp.StatusCode, body)
}
