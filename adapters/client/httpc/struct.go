package httpc

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rendau/dop/adapters/logger"
)

// Options

type OptionsSt struct {
	Client         *http.Client
	Uri            string
	Method         string
	Params         url.Values
	Headers        http.Header
	BasicAuthCreds *BasicAuthCredsSt
	LogFlags       int
	LogPrefix      string
	RetryCount     int
	RetryInterval  time.Duration
	Timeout        time.Duration

	ReqClose     bool
	ReqStream    io.Reader
	ReqBody      []byte
	ReqObj       any
	RepStream    bool
	RepObj       any
	StatusRepObj map[int]any
}

type BasicAuthCredsSt struct {
	Username string
	Password string
}

func (o *OptionsSt) GetMergedWith(val *OptionsSt) *OptionsSt {
	res := &OptionsSt{
		Client:         o.Client,
		Uri:            o.Uri + val.Uri,
		Method:         o.Method,
		Params:         url.Values{},
		Headers:        http.Header{},
		BasicAuthCreds: o.BasicAuthCreds,
		LogFlags:       o.LogFlags,
		LogPrefix:      o.LogPrefix + val.LogPrefix,
		RetryCount:     o.RetryCount,
		RetryInterval:  o.RetryInterval,
		Timeout:        o.Timeout,
	}

	// Client
	if val.Client != nil {
		res.Client = val.Client
	}

	// Method
	if val.Method != "" {
		res.Method = val.Method
	}
	if res.Method == "" {
		res.Method = "GET"
	}

	// Params
	for k, v := range o.Params {
		res.Params[k] = v
	}
	for k, v := range val.Params {
		res.Params[k] = v
	}

	// Headers
	for k, v := range o.Headers {
		res.Headers[k] = v
	}
	for k, v := range val.Headers {
		res.Headers[k] = v
	}

	// BasicAuthCreds
	if val.BasicAuthCreds != nil {
		res.BasicAuthCreds = val.BasicAuthCreds
	}

	// LogFlags
	if val.LogFlags != 0 {
		if val.LogFlags < 0 {
			res.LogFlags = 0
		} else {
			res.LogFlags = val.LogFlags
		}
	}

	// RetryCount
	if val.RetryCount != 0 {
		if val.RetryCount < 0 {
			res.RetryCount = 0
		} else {
			res.RetryCount = val.RetryCount
		}
	}

	// RetryInterval
	if val.RetryInterval != 0 {
		if val.RetryInterval < 0 {
			res.RetryInterval = 0
		} else {
			res.RetryInterval = val.RetryInterval
		}
	}

	// Timeout
	if val.Timeout != 0 {
		if val.Timeout < 0 {
			res.Timeout = 0
		} else {
			res.Timeout = val.Timeout
		}
	}

	// ReqStream
	if val.ReqStream != nil {
		res.ReqStream = val.ReqStream
	}

	// ReqBody
	if val.ReqBody != nil {
		res.ReqBody = val.ReqBody
	}

	// ReqObj
	if val.ReqObj != nil {
		res.ReqObj = val.ReqObj
	}

	// RepStream
	if val.RepStream {
		res.RepStream = val.RepStream
	}

	// RepObj
	if val.RepObj != nil {
		res.RepObj = val.RepObj
	}

	// StatusRepObj
	if val.StatusRepObj != nil {
		res.StatusRepObj = val.StatusRepObj
	}

	return res
}

func (o *OptionsSt) HasLogFlag(v int) bool {
	return o.LogFlags&v > 0
}

// Resp

type RespSt struct {
	Lg      logger.Lite
	ReqOpts *OptionsSt

	StatusCode        int
	StatusCodeSuccess bool
	Headers           http.Header
	BodyRaw           []byte
	Stream            io.ReadCloser
}

func (o *RespSt) Reset() {
	o.StatusCode = 0
	o.StatusCodeSuccess = false
	o.BodyRaw = nil
	o.Stream = nil
}

func (o *RespSt) LogError(title string, err error, args ...any) {
	if o.ReqOpts.HasLogFlag(ErrorLogToInfo) {
		o.LogInfo(title, append(args, "error", err.Error())...)
	} else {
		o.Lg.Errorw(o.ReqOpts.LogPrefix+title, err, o.fillLogArgs(args...)...)
	}
}

func (o *RespSt) LogInfo(title string, args ...any) {
	o.Lg.Infow(o.ReqOpts.LogPrefix+title, o.fillLogArgs(args...)...)
}

func (o *RespSt) fillLogArgs(srcArgs ...any) []any {
	return append(
		srcArgs,
		"method", o.ReqOpts.Method,
		"uri", o.ReqOpts.Uri,
		"params", o.ReqOpts.Params.Encode(),
		"req_body", string(o.ReqOpts.ReqBody),
		"status_code", o.StatusCode,
		"rep_body", string(o.BodyRaw),
	)
}
