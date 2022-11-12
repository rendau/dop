package tests

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/client/httpc/httpclient"
	"github.com/rendau/dop/adapters/server/https"
	"github.com/rendau/dop/dopErrs"
	"github.com/stretchr/testify/require"
)

func TestHttpc(t *testing.T) {
	const ServerPort = "29714"

	sReqObj := struct {
		headers http.Header
		params  url.Values
		body    []byte
	}{}
	var e500Cnt int

	// server

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/s200", func(c *gin.Context) {
		sReqObj.headers = c.Request.Header
		sReqObj.params = c.Request.URL.Query()
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(500)
			return
		}
		sReqObj.body = body
		c.JSON(200, gin.H{"a": "1"})
	})
	r.GET("/e500", func(c *gin.Context) {
		if e500Cnt > 0 {
			e500Cnt--
			c.Status(500)
			return
		}
		c.Status(200)
	})
	r.GET("/e400", func(c *gin.Context) {
		c.JSON(400, gin.H{"error": "code"})
	})
	r.GET("/timeout", func(c *gin.Context) {
		time.Sleep(200 * time.Millisecond)
	})
	r.GET("/basic_auth", gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	}), func(c *gin.Context) {
		c.Status(200)
	})

	server := https.Start(":"+ServerPort, r, app.lg)

	// client

	httpClient := &http.Client{}
	hc := httpclient.New(app.lg, &httpc.OptionsSt{})

	cases := []struct {
		baseOpts *httpc.OptionsSt
		reqOpts  *httpc.OptionsSt
		e500Cnt  int

		wantErr           error
		wantStatusCode    int
		wantSReqHeaders   map[string]string
		wantSReqParams    map[string]string
		wantSReqBodyCheck bool
		wantRepObj        map[string]string
		wantStatusRepObj  map[int]any
	}{
		{
			baseOpts: &httpc.OptionsSt{
				Headers: http.Header{"Authorization": {"token"}},
			},
			reqOpts: &httpc.OptionsSt{
				Uri:      "s200",
				Headers:  http.Header{"Header": {"h_value"}},
				Params:   map[string][]string{"qp": {"qpv"}},
				RepObj:   &map[string]string{},
				LogFlags: httpc.LogRequest | httpc.LogResponse,
			},
			wantErr:         nil,
			wantStatusCode:  200,
			wantSReqHeaders: map[string]string{"Accept": "application/json", "Authorization": "token", "Header": "h_value"},
			wantSReqParams:  map[string]string{"qp": "qpv"},
			wantRepObj:      map[string]string{"a": "1"},
		},
		{
			baseOpts: &httpc.OptionsSt{
				RetryCount:    1,
				RetryInterval: 10 * time.Millisecond,
			},
			reqOpts: &httpc.OptionsSt{
				Uri: "e500",
			},
			e500Cnt:        1,
			wantErr:        nil,
			wantStatusCode: 200,
		},
		{
			baseOpts: &httpc.OptionsSt{
				RetryCount:    1,
				RetryInterval: 10 * time.Millisecond,
			},
			reqOpts: &httpc.OptionsSt{
				Uri: "e500",
			},
			e500Cnt:        3,
			wantErr:        dopErrs.BadStatusCode,
			wantStatusCode: 500,
		},
		{
			baseOpts: &httpc.OptionsSt{},
			reqOpts: &httpc.OptionsSt{
				Uri:     "timeout",
				Timeout: 20 * time.Millisecond,
			},
			wantErr:        context.DeadlineExceeded,
			wantStatusCode: 0,
		},
		{
			baseOpts: &httpc.OptionsSt{},
			reqOpts: &httpc.OptionsSt{
				Uri: "e400",
				StatusRepObj: map[int]any{
					400: &map[string]string{},
				},
			},
			wantErr:        dopErrs.BadStatusCode,
			wantStatusCode: 400,
			wantRepObj:     map[string]string{"a": "1"},
			wantStatusRepObj: map[int]any{
				400: &map[string]string{"error": "code"},
			},
		},
		{
			baseOpts: &httpc.OptionsSt{},
			reqOpts: &httpc.OptionsSt{
				Uri:    "s200",
				ReqObj: map[string]string{"hello": "world"},
				RepObj: &map[string]string{},
			},
			wantStatusCode:    200,
			wantRepObj:        map[string]string{"a": "1"},
			wantSReqBodyCheck: true,
		},
		{
			baseOpts: &httpc.OptionsSt{},
			reqOpts: &httpc.OptionsSt{
				Uri: "basic_auth",
				BasicAuthCreds: &httpc.BasicAuthCredsSt{
					Username: "admin",
					Password: "bad_secret",
				},
			},
			wantErr:        dopErrs.BadStatusCode,
			wantStatusCode: 401,
		},
		{
			baseOpts: &httpc.OptionsSt{
				BasicAuthCreds: &httpc.BasicAuthCredsSt{
					Username: "admin",
					Password: "secret",
				},
			},
			reqOpts: &httpc.OptionsSt{
				Uri: "basic_auth",
			},
			wantStatusCode: 200,
		},
	}

	var err error
	var resp *httpc.RespSt

	for cI, c := range cases {
		t.Run(strconv.Itoa(cI+1), func(t *testing.T) {
			c.baseOpts.Client = httpClient
			c.baseOpts.Uri = "http://localhost:" + ServerPort
			hc.SetOptions(c.baseOpts)

			e500Cnt = c.e500Cnt

			// Send
			resp, err = hc.Send(c.reqOpts)

			// error
			if c.wantErr == nil {
				require.Nil(t, err)
			} else {
				if !errors.Is(err, c.wantErr) {
					require.Equal(t, c.wantErr, err)
				}
			}

			// status-code
			require.Equal(t, c.wantStatusCode, resp.StatusCode)

			// ReqObj
			if c.wantSReqBodyCheck {
				if c.reqOpts.ReqObj != nil {
					raw, err := json.Marshal(c.reqOpts.ReqObj)
					require.Nil(t, err)
					require.JSONEq(t, string(raw), string(sReqObj.body))
				}
			}

			// RepObj && StatusRepObj
			if err == nil {
				if c.wantRepObj != nil {
					require.Equal(t, &c.wantRepObj, c.reqOpts.RepObj)
				}
			} else if errors.Is(err, dopErrs.BadStatusCode) {
				if c.reqOpts.StatusRepObj != nil {
					require.Equal(t, c.wantStatusRepObj, c.reqOpts.StatusRepObj)
				}
			}

			// headers
			for k, v := range c.wantSReqHeaders {
				require.Equal(t, v, sReqObj.headers.Get(k))
			}

			// params
			for k, v := range c.wantSReqParams {
				require.Equal(t, v, sReqObj.params.Get(k))
			}
		})
	}

	require.True(t, server.Shutdown(2*time.Second))
}
