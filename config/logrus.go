package config

import (
	"github.com/evalphobia/logrus_sentry"
	filename "github.com/keepeye/logrus-filename"
	"github.com/sirupsen/logrus"
)

// InitLog 日志服务初始化
func InitLog() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if level, err := logrus.ParseLevel(Cfg.LogLevel); err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	fnhook := filename.NewHook()
	logrus.AddHook(fnhook)

	if len(Cfg.SentryDSN) > 0 {
		if hook, err := logrus_sentry.NewSentryHook(Cfg.SentryDSN, []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.PanicLevel,
		}); err == nil {
			hook.StacktraceConfiguration.Enable = true
			logrus.AddHook(hook)
		} else {
			logrus.WithError(err).Warn("init sentry failed")
		}
	}
}
