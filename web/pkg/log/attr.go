package log

type HTTP struct {
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
}

type Error struct {
	Message string `json:"message"`
}

type Network struct {
	ClientIP string `json:"client_ip"`
}

type DataDog struct {
	TraceID uint64 `json:"trace_id"`
	SpanID  uint64 `json:"span_id"`
}
