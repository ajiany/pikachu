package config

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

// StartDDTrace datadog服务初始化
func StartDDTrace() {
	tracer.Start()
	logrus.Printf("dd trace started")
	if Cfg.DDTrace {
		if err := profiler.Start(); err != nil {
			logrus.WithError(err).Warn("start dd trace profiler failed")
		}
		logrus.Printf("dd profiler started")
	}
}

func StopDDTrace(timeout time.Duration) {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		tracer.Stop()
		logrus.Printf("dd trace stop")
		return nil
	})
	if Cfg.DDTrace {
		eg.Go(func() error {
			profiler.Stop()
			logrus.Printf("dd profiler stop")
			return nil
		})
	}
	eg.Wait()
}
