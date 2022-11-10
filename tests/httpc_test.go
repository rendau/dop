package tests

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rendau/dop/adapters/client/httpc"
	"github.com/rendau/dop/adapters/client/httpc/httpclient"
	"github.com/rendau/dop/adapters/server/https"
	"github.com/stretchr/testify/require"
)

func TestHttpc(t *testing.T) {
	const ServerPort = "29714"

	var err error

	sReqObj := struct {
		headers http.Header
		params  url.Values
	}{}
	var e500Cnt int

	// server

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/s200", func(c *gin.Context) {
		sReqObj.headers = c.Request.Header
		sReqObj.params = c.Request.URL.Query()
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

	server := https.Start(":"+ServerPort, r, app.lg)

	// client

	httpClient := &http.Client{}
	hc := httpclient.New(app.lg, httpc.OptionsSt{})

	cases := []struct {
		baseOpts httpc.OptionsSt
		reqOpts  httpc.OptionsSt
		e500Cnt  int

		wantErr         error
		wantStatusCode  int
		wantSReqHeaders map[string]string
		wantSReqParams  map[string]string
		wantRepObj      map[string]string
	}{
		{
			baseOpts: httpc.OptionsSt{
				Headers: http.Header{"Authorization": {"token"}},
			},
			reqOpts: httpc.OptionsSt{
				Uri:      "s200",
				Headers:  http.Header{"Header": {"h_value"}},
				Params:   map[string][]string{"qp": {"qpv"}},
				LogFlags: httpc.LogRequest | httpc.LogResponse,
			},
			wantErr:         nil,
			wantStatusCode:  200,
			wantSReqHeaders: map[string]string{"Accept": "application/json", "Authorization": "token", "Header": "h_value"},
			wantSReqParams:  map[string]string{"qp": "qpv"},
			wantRepObj:      map[string]string{"a": "1"},
		},
		{
			baseOpts: httpc.OptionsSt{
				RetryCount:    1,
				RetryInterval: 10 * time.Millisecond,
			},
			reqOpts: httpc.OptionsSt{
				Uri: "e500",
			},
			e500Cnt:        1,
			wantErr:        nil,
			wantStatusCode: 200,
		},
	}

	var resp *httpc.RespSt

	for cI, c := range cases {
		t.Run(strconv.Itoa(cI+1), func(t *testing.T) {
			c.baseOpts.Client = httpClient
			c.baseOpts.Uri = "http://localhost:" + ServerPort
			hc.SetOptions(c.baseOpts)

			e500Cnt = c.e500Cnt

			if c.reqOpts.ReqObj == nil {
				if c.wantRepObj == nil {
					resp, err = hc.Send(c.reqOpts)
					require.Equal(t, c.wantErr, err)
					require.Equal(t, c.wantStatusCode, resp.StatusCode)
				} else {
					c.reqOpts.RepObj = &map[string]string{}
					resp, err = hc.SendRecvJson(c.reqOpts)
					require.Equal(t, c.wantErr, err)
					require.Equal(t, c.wantStatusCode, resp.StatusCode)
					require.Equal(t, &c.wantRepObj, c.reqOpts.RepObj)
				}
			} else {
				if c.wantRepObj == nil {
				} else {
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
