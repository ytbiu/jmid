package jmid

type logCtxKey struct{}

var (
	TraceLogMicroCtxKey = logCtxKey{}
)

const (
	traceIDKey  = "x-atrace-id"
	traceLogKey = "trace-log"
)
