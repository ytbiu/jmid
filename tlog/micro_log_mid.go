package tlog

import (
	"jmid"

	"context"

	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/sirupsen/logrus"
	"gitlab.exmarttech.com/exsmart-go/warpper/mid"

)

// micro opentracing中间件
func MicroLogMid() server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			md, exist := metadata.FromContext(ctx)
			if !exist {
				h(ctx, req, rsp)
				return nil
			}
			traceID, exist := md[jmid.TraceIDKey]
			if !exist {
				h(ctx, req, rsp)
				return nil
			}

			logWithTraceID := logrus.WithField(mid.TraceIDKey, traceID)
			ctx = context.WithValue(ctx, mid.TraceLogMicroCtxKey, logWithTraceID)

			h(ctx, req, rsp)
			return nil
		}
	}
}

// 获取注入traceid后的日志对象
func FromCtx(c context.Context) *logrus.Entry {
	logWithTrace, exist := c.Value(mid.TraceLogMicroCtxKey).(*logrus.Entry)
	if !exist {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return logWithTrace
}
