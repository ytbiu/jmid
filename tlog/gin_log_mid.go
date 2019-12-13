package tlog

import (
	"jmid"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// gin opentracing 中间件
func GinTraceMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if traceID := c.Request.Header.Get(jmid.TraceIDKey); traceID != "" {
			logWithTraceID := logrus.WithField(jmid.TraceIDKey, c.Request.Header.Get(jmid.TraceIDKey))
			c.Set(jmid.TraceLogKey, logWithTraceID)
		}
	}

}

// gin 获取注入traceid 后的日志对象
func FromGCtx(c *gin.Context) *logrus.Entry {
	v, exist := c.Get(jmid.TraceLogKey)
	if !exist {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return v.(*logrus.Entry)
}
