package main

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Infof("CONTEXT: %v", ctx)
	log.WithFields(
		log.Fields{
			"subdomain": ctx.Value("subdomain"),
			"trace_id":  ctx.Value("datadogTraceID"),
		},
	).Errorf(format, args...)
}
