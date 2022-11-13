package form

import (
	"context"
	"github.com/ajiany/pikachu/tools/test_helper"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/ajiany/pikachu/tools/test_helper"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	message = "oh yeah"
	code    = 429
	status  = 1
)

type customForm struct{}

func (form customForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewCustomError(code, status, "")
}

type successForm struct{}

func (form successForm) Do(c *gin.Context) (interface{}, error) {
	return nil, nil
}

type notFoundForm struct{}

func (form notFoundForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewNotFoundError(message)
}

type unprocessableForm struct{}

func (form unprocessableForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewUnprocessableError(message)
}

type internalStateForm struct{}

func (form internalStateForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewInternalStateError(message)
}

type unauthorizedForm struct{}

func (form unauthorizedForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewUnauthorizedError(message)
}

type withCodeForm struct{}

func (form withCodeForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewWithCodeError(code, message)
}

type withStatusForm struct{}

func (form withStatusForm) Do(c *gin.Context) (interface{}, error) {
	return nil, NewWithStatusError(status, message)
}

type cancelForm struct{}

func (form cancelForm) Do(c *gin.Context) (interface{}, error) {
	return nil, context.Canceled
}

func TestFormSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cases := []struct {
		form                 Form
		customCode           bool
		status               int
		withStatusInResponse bool
		statusInResponse     int
		message              string
	}{
		{&customForm{}, true, 0, true, status, ""},
		{&successForm{}, false, 200, false, 0, ""},
		{&notFoundForm{}, false, 404, false, 0, message},
		{&unauthorizedForm{}, false, 401, false, 0, message},
		{&unprocessableForm{}, false, 422, false, 0, message},
		{&internalStateForm{}, false, 520, false, 0, message},
		{&withCodeForm{}, true, 0, false, 0, message},
		{&withStatusForm{}, false, 200, true, status, message},
		{&cancelForm{}, false, 408, false, 0, "context canceled"},
	}

	for _, tc := range cases {
		r := gin.Default()
		r.GET("/test", func(c *gin.Context) {
			ok, _ := DoThisForm(tc.form, c)
			if ok {
				c.String(200, "")
			}
		})

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if tc.customCode {
			assert.Equal(t, w.Code, code)
		} else {
			assert.Equal(t, tc.status, w.Code)
		}

		if tc.withStatusInResponse {
			js := test_helper.JSONResp(w.Result())
			assert.Equal(t, tc.message, js.Get("msg").MustString())
			assert.Equal(t, tc.statusInResponse, js.Get("status").MustInt())
		} else {
			if w.Code != 200 {
				js := test_helper.JSONResp(w.Result())
				assert.Equal(t, tc.message, js.Get("errors").GetIndex(0).MustString())
			}
		}
	}
}
