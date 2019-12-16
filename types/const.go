package types

type logCtxKey struct{}

var (
	TraceLogMicroCtxKey = logCtxKey{}
)

const (
	TraceIDKey  = "x-atrace-id"
	TraceLogKey = "trace-log"
)
