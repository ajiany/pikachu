package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

var defaultLogHeaders = map[string]string{
	"request_id":   "Request-Id",
	"origin":       "Origin",
	"locale":       "X-LOCALE",
	"referer":      "referer",
	"x_request_id": "x-request-id",
}

type requestInfo struct {
	ClientIp  string
	Status    int
	Elapse    time.Duration
	Path      string
	Query     string
	Method    string
	UserAgent string

	customHeaders map[string]string
	headers       http.Header
}

func (s *requestInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("client_ip", s.ClientIp)
	enc.AddInt("status", s.Status)
	enc.AddDuration("elapse", s.Elapse)
	enc.AddString("path", s.Path)
	enc.AddString("query", s.Query)
	enc.AddString("method", s.Method)
	enc.AddString("user_agent", s.UserAgent)

	for k, v := range defaultLogHeaders {
		h := s.headers.Get(v)
		if len(h) > 0 {
			enc.AddString(k, h)
		}
	}

	if len(s.customHeaders) > 0 {
		for k, v := range s.customHeaders {
			h := s.headers.Get(v)
			if len(h) > 0 {
				enc.AddString(k, h)
			}
		}
	}
	return nil
}

func AccessLog(service string, setters ...AccessLogOption) gin.HandlerFunc {
	cfg := zap.NewProductionConfig()
	cfg.InitialFields = map[string]interface{}{"service": service}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	opts := &AccessLogOptions{}
	for _, setter := range setters {
		setter(opts)
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		logger.Info(path,
			zap.String("log_type", "access"),
			zap.String("time", time.Now().Format(time.RFC3339)),

			zap.Object("message", &requestInfo{
				Status:        c.Writer.Status(),
				Method:        c.Request.Method,
				Path:          path,
				Query:         query,
				ClientIp:      c.ClientIP(),
				UserAgent:     c.Request.UserAgent(),
				Elapse:        latency,
				customHeaders: opts.CustomHeaders,
				headers:       c.Request.Header,
			}),
		)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors.Errors() {
				logger.Error(e)
			}
		}
	}
}

type AccessLogOption func(*AccessLogOptions)

type AccessLogOptions struct {
	CustomHeaders map[string]string
}

// AccessLogOptHeaders
// customHeaders Example: map[string]string{"referer": "referer", "x_request_id": "x-request-id"}
// key: log message key
// value: the header that you want to add to the log
func AccessLogOptHeaders(customHeaders map[string]string) AccessLogOption {
	return func(args *AccessLogOptions) {
		args.CustomHeaders = customHeaders
	}
}
