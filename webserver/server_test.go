package webserver

import (
	"context"
	testhelper "github.com/ajiany/pikachu/tools/test_helper"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
	"time"
)

func TestStartAndClose(t *testing.T) {

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	server, err := StartHTTPServer(context.Background(), r)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	resp := testhelper.RunGetRequest(server.URL() + "/ping")

	Convey("http server is started", t, func() {
		So(resp.StatusCode, ShouldEqual, 200)
		So(testhelper.StringResp(resp), ShouldEqual, "pong")
	})

}

func TestStartAndCancel(t *testing.T) {

	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	server, err := StartHTTPServer(ctx, r)
	if err != nil {
		panic(err)
	}

	resp := testhelper.RunGetRequest(server.URL() + "/ping")

	Convey("http server is started", t, func() {
		So(resp.StatusCode, ShouldEqual, 200)
		So(testhelper.StringResp(resp), ShouldEqual, "pong")
	})

	server.Block()
}
