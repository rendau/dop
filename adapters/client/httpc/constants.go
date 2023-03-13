package httpc

const (
	LogRequest            = 1
	LogResponse           = 2
	NoLogError            = 4
	ErrorLogToInfo        = 8
	NoLogNotAuthorized    = 16
	NoLogPermissionDenied = 32
	NoLogBadStatus        = 64
)
