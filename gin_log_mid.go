package jmid

import (

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// gin opentracing 中间件
func TraceGinLogMid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if traceID := c.Request.Header.Get(traceIDKey); traceID != "" {
			logWithTraceID := logrus.WithField(traceIDKey, c.Request.Header.Get(traceIDKey))
			c.Set(traceLogKey, logWithTraceID)
		}
	}

}

// gin 获取注入traceid 后的日志对象
func FromGCtx(c *gin.Context) *logrus.Entry {
	v, exist := c.Get(traceLogKey)
	if !exist {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return v.(*logrus.Entry)
}
