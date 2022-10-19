package feishu

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidSign  = errors.New("Invalid Sign")
	ErrInvalidToken = errors.New("Invalid Token")
)

type ErrResponse struct {
	ResponseData
}

func NewErrResponse(data ResponseData) ErrResponse {
	return ErrResponse{
		ResponseData: data,
	}
}

func (err ErrResponse) Error() string {
	return fmt.Sprintf("Request fail: code -> %d, msg -> %s", err.Code, err.Msg)
}
