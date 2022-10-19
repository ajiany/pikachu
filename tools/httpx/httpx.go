package httpx

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

type Httpx struct {
	logger *logrus.Logger
	cfg    *Config

	beforeBuildRequest BeforeRequestHook
	beforeRequest      RequestHook
	afterResponse      ResponseHook
}

func NewHttpx(fns ...configFn) *Httpx {
	cfg := defaultCfg()
	for _, fn := range fns {
		fn(cfg)
	}

	return &Httpx{
		logger: logrus.StandardLogger(),
		cfg:    cfg,
	}
}

func GET(ctx context.Context, url string, opt *Option) (*Response, error) {
	return NewHttpx().GET(ctx, url, opt)
}

func POST(ctx context.Context, url string, opt *Option) (*Response, error) {
	return NewHttpx().POST(ctx, url, opt)
}

func PUT(ctx context.Context, url string, opt *Option) (*Response, error) {
	return NewHttpx().PUT(ctx, url, opt)
}

func PATCH(ctx context.Context, url string, opt *Option) (*Response, error) {
	return NewHttpx().PATCH(ctx, url, opt)
}

func DELETE(ctx context.Context, url string, opt *Option) (*Response, error) {
	return NewHttpx().DELETE(ctx, url, opt)
}

func (h *Httpx) GET(ctx context.Context, url string, opt *Option) (*Response, error) {
	return h.Request(ctx, http.MethodGet, url, opt)
}

func (h *Httpx) DoGET(ctx context.Context, url string, opt *Option) (*Response, error) {
	return h.Request(ctx, http.MethodGet, url, opt)
}

func (h *Httpx) POST(ctx context.Context, url string, opt *Option) (*Response, error) {
	return h.Request(ctx, http.MethodPost, url, opt)
}

func (h *Httpx) PUT(ctx context.Context, url string, opt *Option) (*Response, error) {
	return h.Request(ctx, http.MethodPut, url, opt)
}

func (h *Httpx) PATCH(ctx context.Context, url string, opt *Option) (*Response, error) {
	return h.Request(ctx, http.MethodPatch, url, opt)
}

func (h *Httpx) DELETE(ctx context.Context, url string, opt *Option) (*Response, error) {
	return h.Request(ctx, http.MethodDelete, url, opt)
}

func (h *Httpx) Request(ctx context.Context, method string, url string, opt *Option) (*Response, error) {
	req, err := h.buildRequest(ctx, method, url, opt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cli := http.DefaultClient
	if h.cfg.Trace {
		cli = httptrace.WrapClient(cli)
	}
	if h.cfg.Timeout != 0 {
		cli.Timeout = h.cfg.Timeout
	}

	var code int
	start := time.Now().UnixNano()
	defer func() {
		duration := float64(time.Now().UnixNano()-start) / 1000000.0
		h.logger.Infof("%d <--- %.6f ms", code, duration)
	}()

	if opt != nil && opt.Params != nil {
		b, _ := opt.Parsed()
		h.logger.Infof("[%s] %s %+v", method, req.Req.URL.String(), string(b))
	} else {
		h.logger.Infof("[%s] %s", method, req.Req.URL.String())
	}

	h.logger.Debugf("request header -> %v", req.Req.Header)

	resp, err := cli.Do(req.Req)
	r := newResponse(resp)
	if err != nil {
		return r, errors.Wrap(err, "do request fail")
	}

	code = r.Resp.StatusCode
	respBody, _ := r.ParsedString()
	h.logger.Debugf("%d -> %s", r.Resp.StatusCode, Escape(Truncation(respBody, 200)))
	h.logger.Debugf("response header -> %v", r.Resp.Header)

	if !r.Success() {
		return r, NewRequestFailErr(r)
	}

	if !opt.SkipHooks {
		if err := h.runResponseHook(r); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (h *Httpx) buildRequest(ctx context.Context, method string, url string, opt *Option) (*Request, error) {
	reqUrl := h.getRequestUrl(url, opt)

	if !opt.SkipHooks {
		h.runBeforeBuildRequestHook(opt)
	}

	req, err := newRequest(ctx, method, reqUrl, opt)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !opt.SkipHooks {
		h.runBeforeHook(req, opt)
	}

	// headers
	if len(opt.Headers) != 0 {
		for k, v := range opt.Headers {
			for _, item := range v {
				req.Req.Header.Add(k, item)
			}
		}
	}

	return req, nil
}

func (h *Httpx) getRequestUrl(url string, opt *Option) string {
	if h.cfg.Host == "" || strings.HasPrefix(url, "http") {
		return url
	}

	if opt.SkipUrlPrefix {
		return fmt.Sprintf("%s%s", h.cfg.Host, url)
	}

	return fmt.Sprintf("%s%s%s", h.cfg.Host, h.cfg.UrlPrefix, url)
}

func (h *Httpx) SetBeforeBuildRequestHook(fn func(opts *Option)) {
	h.beforeBuildRequest = fn
}

func (h *Httpx) SetBeforeRequestHook(fn func(req *Request, opts *Option)) {
	h.beforeRequest = fn
}

func (h *Httpx) SetAfterRequestHook(fn ResponseHook) {
	h.afterResponse = fn
}

func (h *Httpx) runBeforeBuildRequestHook(opts *Option) {
	if h.beforeBuildRequest != nil {
		h.beforeBuildRequest(opts)
	}
}

func (h *Httpx) runBeforeHook(req *Request, opts *Option) {
	if h.beforeRequest != nil {
		h.beforeRequest(req, opts)
	}
}

func (h *Httpx) runResponseHook(resp *Response) error {
	if h.afterResponse == nil {
		return nil
	}

	if err := h.afterResponse(resp); err != nil {
		return err
	}

	return nil
}
