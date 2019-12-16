package mjaeger

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"jmid/types"
)

const (
	TRACE_CONTEXT_KEY = "trace-context"
)

// gin jaeger中间件
func GinMid(tags ...map[string]interface{}) gin.HandlerFunc {

	return func(c *gin.Context) {

		var span opentracing.Span
		defer func() {
			// set tags
			statusCode := c.Writer.Status()
			ext.HTTPStatusCode.Set(span, uint16(statusCode))
			ext.HTTPMethod.Set(span, c.Request.Method)
			ext.HTTPUrl.Set(span, c.Request.URL.EscapedPath())
			if statusCode >= http.StatusInternalServerError {
				ext.Error.Set(span, true)
			}

			// set extra tags
			if len(tags) > 0 {
				for k, v := range tags[0] {
					span.SetTag(k, v)
				}
			}
			span.Finish()
		}()

		// get a span
		md := make(map[string]string)
		spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		if err != nil {
			span = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path)
		} else {
			span = opentracing.GlobalTracer().StartSpan(c.Request.URL.Path, opentracing.ChildOf(spanCtx))
		}

		traceID := c.Request.Header.Get(types.TraceIDKey)
		md[types.TraceIDKey] = traceID
		span.SetTag(types.TraceIDKey, traceID)

		if err := opentracing.GlobalTracer().Inject(span.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(md)); err != nil {
			log.Println(err)
		}

		// span注入context
		ctx := opentracing.ContextWithSpan(context.TODO(), span)
		ctx = metadata.NewContext(ctx, md)
		// context,span 注入gin context
		c.Set(TRACE_CONTEXT_KEY, ctx)

		c.Next()
	}
}

// ContextWithSpan 返回context
func ContextWithSpan(c *gin.Context) (ctx context.Context) {
	v, exist := c.Get(TRACE_CONTEXT_KEY)
	if !exist {
		return context.TODO()
	}

	ctx, exist = v.(context.Context)
	if !exist {
		return context.TODO()
	}

	return ctx
}
