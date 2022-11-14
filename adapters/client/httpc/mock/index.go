package mock

import (
	"encoding/json"
	"sync"

	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/logger"
	"github.com/rendau/dop/dopErrs"
)

const (
	ErrPageNotFound = dopErrs.Err("page_not_found")
)

type St struct {
	lg logger.Lite

	requests  []*httpc.OptionsSt
	responses map[string]ResponseSt
	mu        sync.Mutex
}

type ResponseSt struct {
	RespObj any
	Resp    *httpc.RespSt
	Err     error
}

func New(lg logger.Lite) *St {
	return &St{
		lg: lg,

		requests:  []*httpc.OptionsSt{},
		responses: map[string]ResponseSt{},
	}
}

func (c *St) SetResponses(responses map[string]ResponseSt) {
	c.mu.Lock()
	c.responses = map[string]ResponseSt{}
	c.mu.Unlock()

	for k, v := range responses {
		c.SetResponse(k, v)
	}
}

func (c *St) SetResponse(path string, response ResponseSt) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if response.Resp == nil {
		response.Resp = &httpc.RespSt{}
	}

	response.Resp.Lg = c.lg

	if len(response.Resp.BodyRaw) == 0 && response.RespObj != nil {
		var err error

		response.Resp.BodyRaw, err = json.Marshal(response.RespObj)
		if err != nil {
			c.lg.Errorw("Fail to marshal json", err)
		}
	}

	c.responses[path] = response
}

func (c *St) SetOptions(opts *httpc.OptionsSt) {
}

func (c *St) GetOptions() *httpc.OptionsSt {
	return &httpc.OptionsSt{}
}

func (c *St) Send(opts *httpc.OptionsSt) (*httpc.RespSt, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error

	if opts.ReqObj != nil {
		opts.ReqBody, err = json.Marshal(opts.ReqObj)
		if err != nil {
			c.lg.Errorw("Fail to marshal json", err)
			return nil, err
		}
	}

	c.requests = append(c.requests, opts)

	response, ok := c.responses[opts.Uri]
	if !ok {
		c.lg.Infow("Httpc-mock, path not found", "path", opts.Uri)
		return nil, ErrPageNotFound
	}

	if len(response.Resp.BodyRaw) > 0 && opts.RepObj != nil {
		err = json.Unmarshal(response.Resp.BodyRaw, opts.RepObj)
		if err != nil {
			c.lg.Errorw("Fail to unmarshal json", err)
			return nil, err
		}
	}

	return response.Resp, nil
}

func (c *St) GetRequests() []*httpc.OptionsSt {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make([]*httpc.OptionsSt, len(c.requests))

	for i, req := range c.requests {
		result[i] = req
	}

	return result
}

func (c *St) GetRequest(path string, obj any) (*httpc.OptionsSt, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, req := range c.requests {
		if req.Uri != path {
			continue
		}

		if len(req.ReqBody) > 0 && obj != nil {
			err := json.Unmarshal(req.ReqBody, obj)
			if err != nil {
				c.lg.Errorw("Fail to unmarshal json", err)
				return nil, false
			}
		}

		return req, true
	}

	return nil, false
}

func (c *St) Clean() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.requests = []*httpc.OptionsSt{}
	c.responses = map[string]ResponseSt{}
}
