package jmid

type logCtxKey struct{}

var (
	traceLogMicroCtxKey = logCtxKey{}
)

const (
	traceIDKey  = "x-atrace-id"
	traceLogKey = "trace-log"
)
