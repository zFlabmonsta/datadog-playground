package log

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Fields = log.Fields
type Formatter = log.Formatter

type LoggerWrapper struct {
	*log.Logger
}

func NewLogger(l *log.Logger, f Formatter) *LoggerWrapper {
	logger := &LoggerWrapper{
		l,
	}
	l.Formatter = f
	return logger
}

func (l *LoggerWrapper) Errorf(ctx context.Context, format string, args ...interface{}) {
	l.Logger.WithFields(l.processAttrFields(ctx)).Errorf(format, args...)
}

func (l *LoggerWrapper) Infof(ctx context.Context, format string, args ...interface{}) {
	l.Logger.WithFields(l.processAttrFields(ctx)).Infof(format, args...)
}

func (l *LoggerWrapper) processAttrFields(ctx context.Context) Fields {
	return Fields{
		"subdomain": ctx.Value("subdomain"),
		"dd": DataDog{
			TraceID: ctx.Value("datadogTraceID").(uint64),
			SpanID:  ctx.Value("datadogSpanID").(uint64),
		},
		"http": ctx.Value(HTTP{}),
	}
}
