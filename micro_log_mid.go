package jmid

import (
	"context"

	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/sirupsen/logrus"
)

// micro opentracing中间件
func TraceMicroLogMid() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			md, exist := metadata.FromContext(ctx)
			if !exist {
				h(ctx, req, rsp)
				return nil
			}
			traceID, exist := md[traceIDKey]
			if !exist {
				h(ctx, req, rsp)
				return nil
			}

			logWithTraceID := logrus.WithField(traceIDKey, traceID)
			ctx = context.WithValue(ctx, traceLogMicroCtxKey, logWithTraceID)

			h(ctx, req, rsp)
			return nil
		}
	}
}

// 获取注入traceid后的日志对象
func FromCtx(c context.Context) *logrus.Entry {
	logWithTrace, exist := c.Value(traceLogMicroCtxKey).(*logrus.Entry)
	if !exist {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logWithTrace
}
