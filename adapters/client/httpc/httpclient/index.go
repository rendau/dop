package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/dop/dopErrs"
)

type St struct {
	lg   logger.Lite
	opts *httpc.OptionsSt
}

func New(lg logger.Lite, opts *httpc.OptionsSt) *St {
	res := &St{
		lg: lg,
	}

	res.SetOptions(opts)

	return res
}

func (c *St) GetOptions() *httpc.OptionsSt {
	return c.opts
}

func (c *St) SetOptions(opts *httpc.OptionsSt) {
	if opts.Uri != "" {
		opts.Uri = strings.TrimRight(opts.Uri, "/") + "/"
	}

	c.opts = opts
}

func (c *St) Send(opts *httpc.OptionsSt) (*httpc.RespSt, error) {
	var err error

	opts = c.opts.GetMergedWith(opts)

	resp := &httpc.RespSt{ReqOpts: opts, Lg: c.lg}

	// ReqObj
	if len(opts.Headers.Values("Content-Type")) == 0 {
		opts.Headers["Content-Type"] = []string{"application/json"}
	}
	opts.ReqBody, err = json.Marshal(opts.ReqObj)
	if err != nil {
		c.lg.Errorw(opts.LogPrefix+"Fail to marshal json", err)
		return resp, err
	}

	// RepObj
	if opts.RepObj != nil {
		if len(opts.Headers.Values("Accept")) == 0 {
			opts.Headers["Accept"] = []string{"application/json"}
		}
	}

	if opts.HasLogFlag(httpc.LogRequest) {
		resp.LogInfo("Request: " + opts.Uri)
	}

	for i := opts.RetryCount; i >= 0; i-- {
		resp.Reset()
		err = c.send(opts, resp)
		if err == nil {
			if resp.StatusCode < 500 { // not retry for "< 500" errors
				break
			}
		}
		if opts.RetryInterval > 0 && i > 0 {
			time.Sleep(opts.RetryInterval)
		}
	}

	if err == nil {
		err = c.handleRespBadStatusCode(resp)
		if err == nil {
			// RepObj
			if len(resp.BodyRaw) > 0 && resp.ReqOpts.RepObj != nil {
				err = json.Unmarshal(resp.BodyRaw, resp.ReqOpts.RepObj)
				if err != nil {
					if !resp.ReqOpts.HasLogFlag(httpc.NoLogError) {
						resp.LogError("Fail to unmarshal body", err)
					}
					return resp, err
				}
			}

			if opts.HasLogFlag(httpc.LogResponse) {
				resp.LogInfo("Response: " + opts.Uri)
			}
		}
	} else {
		if !opts.HasLogFlag(httpc.NoLogError) {
			resp.LogError("Fail to send http-request", err)
		}
	}

	return resp, err
}

func (c *St) send(opts *httpc.OptionsSt, resp *httpc.RespSt) error {
	var err error

	var req *http.Request

	if opts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()
		req, err = http.NewRequestWithContext(ctx, opts.Method, opts.Uri, bytes.NewBuffer(opts.ReqBody))
	} else {
		req, err = http.NewRequest(opts.Method, opts.Uri, bytes.NewBuffer(opts.ReqBody))
	}
	if err != nil {
		return err
	}

	// headers
	req.Header = opts.Headers

	// params
	req.URL.RawQuery = opts.Params.Encode()

	// Basic auth
	if opts.BasicAuthCreds != nil {
		req.SetBasicAuth(opts.BasicAuthCreds.Username, opts.BasicAuthCreds.Password)
	}

	// Do request
	rep, err := opts.Client.Do(req)
	if err != nil {
		return err
	}
	defer rep.Body.Close()

	// read response body
	resp.BodyRaw, err = ioutil.ReadAll(rep.Body)
	if err != nil {
		return err
	}

	resp.StatusCode = rep.StatusCode

	return nil
}

func (c *St) handleRespBadStatusCode(resp *httpc.RespSt) error {
	if resp.StatusCode > 0 && (resp.StatusCode < 200 || resp.StatusCode > 299) {
		if sObj, ok := resp.ReqOpts.StatusRepObj[resp.StatusCode]; ok {
			if len(resp.BodyRaw) > 0 {
				err := json.Unmarshal(resp.BodyRaw, sObj)
				if err != nil {
					if !resp.ReqOpts.HasLogFlag(httpc.NoLogError) {
						resp.LogError("Fail to unmarshal body", err)
					}
					return err
				}
			}
		} else if !resp.ReqOpts.HasLogFlag(httpc.NoLogError) && !resp.ReqOpts.HasLogFlag(httpc.NoLogBadStatus) {
			switch {
			case resp.StatusCode == 401 && resp.ReqOpts.HasLogFlag(httpc.NoLogNotAuthorized):
			case resp.StatusCode == 403 && resp.ReqOpts.HasLogFlag(httpc.NoLogPermissionDenied):
			default:
				resp.LogError("Bad status code", dopErrs.BadStatusCode)
			}
		}
		return dopErrs.BadStatusCode
	}

	return nil
}
